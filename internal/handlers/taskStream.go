package handlers

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Data to be broadcasted to a client.
type Data struct {
	Message string `json:"message"`
	To      uint   `json:"receiver"`
	Type    string `json:"type"`
}

// Uniquely defines an incoming client.
type Client struct {
	ID      uint
	Channel chan Data
}

// Keeps track of every SSE events.
type Stream struct {
	Message chan Data

	// New client connections
	NewClients chan Client

	// Closed client connections

	ClosedClients chan Client

	//Client to Channel map

	ClientChannels map[uint][]chan Data
}

var lock = &sync.Mutex{}

var stream *Stream

func GetTaskStream() (stream *Stream) {
	if stream == nil {
		lock.Lock()
		defer lock.Unlock()
		if stream == nil {
			stream = &Stream{
				Message:        make(chan Data, 50),
				NewClients:     make(chan Client, 50),
				ClosedClients:  make(chan Client, 50),
				ClientChannels: make(map[uint][]chan Data),
			}
			go stream.listen()
		}
	}
	return
}

func (stream *Stream) listen() {
	for {
		select {
		// Add new client connection
		case client := <-stream.NewClients:
			stream.ClientChannels[client.ID] = append(stream.ClientChannels[client.ID], client.Channel)
			log.Printf("Added client. %d registered clients", len(stream.ClientChannels))
		// Remove closed client connection
		case client := <-stream.ClosedClients:
			channels := stream.ClientChannels[client.ID]
			for i, ch := range channels {
				//remove channel from slice
				if ch == client.Channel {
					stream.ClientChannels[client.ID] = append(channels[:i], channels[i+1:]...)
					break
				}
			}
			close(client.Channel)
			if len(stream.ClientChannels[client.ID]) == 0 {
				delete(stream.ClientChannels, client.ID)
				log.Printf("Removed client. %d registered clients", len(stream.ClientChannels))
			}
		// Broadcast message to a specific client with client ID fetched from eventMsg.To
		case eventMsg := <-stream.Message:
			channels := stream.ClientChannels[eventMsg.To]
			if len(channels) == 0 {
				log.Printf("Receiver - %d not connected", eventMsg.To)
				continue
			}

			for _, ch := range channels {
				select {
				case ch <- eventMsg:
				default:
					log.Printf("Dropping event for user %d (slow tab)", eventMsg.To)
				}
			}
		}

	}
}

func (stream *Stream) SSEConnMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		// Increment global variable ID
		ID := gctx.GetUint("user_id")
		// Initialize client
		client := Client{
			ID:      ID,
			Channel: make(chan Data, 10),
		}

		// Send new connection to event to store
		stream.NewClients <- client

		defer func() {
			// Send closed connection to event server
			log.Printf("Closing connection : %d", client.ID)
			stream.ClosedClients <- client
		}()
		gctx.Writer.Header().Set("Content-Type", "text/event-stream")
		gctx.Writer.Header().Set("Cache-Control", "no-cache")
		gctx.Writer.Header().Set("Connection", "keep-alive")
		gctx.Writer.Header().Set("Transfer-Encoding", "chunked")
		gctx.Set("client", client)
		gctx.Next()
	}
}

func (stream *Stream) SendEvent(event Data) {
	stream.Message <- event

}

func (stream *Stream) TaskStreamHandler(c *gin.Context) {
	v, ok := c.Get("client")
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	client, ok := v.(Client)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-client.Channel:
			if !ok {
				return false
			}
			c.SSEvent("message", msg)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}
