package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/config"
	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/schemas"
	"github.com/syed.fazil/vtask/internal/services"
	"gorm.io/gorm"
)

func SubscribePush(c *gin.Context, db *gorm.DB) {
	var input schemas.SubscribeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	strUserId, exists := c.Get("user_id")

	if !exists {
		log.Println("User id not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	userId, ok := strUserId.(uint)
	if !ok {
		log.Println("Invalid user")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	sub := models.PushSubscription{
		UserID:   userId,
		Endpoint: input.Endpoint,
		P256dh:   input.Keys.P256dh,
		Auth:     input.Keys.Auth,
	}

	if err := db.WithContext(c.Request.Context()).Where("endpoint = ?", sub.Endpoint).Assign(sub).FirstOrCreate(&sub).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save subscription"})
		return
	}

	c.JSON(200, gin.H{"message": "Subscribed"})
}

func TestPushHandler(c *gin.Context, db *gorm.DB) {
	var subs []models.PushSubscription
	db.Find(&subs)
	pushService := services.NewPushService(c.Request.Context(), *config.App)
	payload := services.Payload{
		Title:   "Test Push 🎉",
		Body:    "Sent from Go backend",
		Vibrate: true,
	}
	payload.Data.URL = "/"

	for _, sub := range subs {
		pushService.Send(sub, payload)
	}
	c.JSON(200, gin.H{"message": "Push sent"})

}
