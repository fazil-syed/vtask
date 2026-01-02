package parser

import (
	"context"
	"log"

	pb "github.com/syed.fazil/vtask/internal/proto/nlp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Intent struct {
	Name     string // e.g., "create_task", "mark_done", "query"
	Title    string // extracted task title (if any)
	Reminder string // RFC3339 string or empty
}

func ParseIntent(ctx context.Context, transcription string) Intent {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to nlp service: %v", err)
	}
	defer conn.Close()
	nlpClient := pb.NewNLPServiceClient(conn)
	response, err := nlpClient.Extract(ctx, &pb.ExtractRequest{Text: transcription})
	if err != nil {
		log.Fatalf("failed to get the parsed intent: %v", err)
	}
	task := response.Task
	intent := response.Intent
	iso_time := response.TimeIso

	return Intent{
		Name:     intent,
		Title:    task,
		Reminder: iso_time,
	}
}
