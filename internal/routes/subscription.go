package routes

// add routes for handling tasks such as creating a task, marking a task as done, querying tasks
import (
	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/handlers"
	"github.com/syed.fazil/vtask/internal/middlewares"
	"gorm.io/gorm"
)

func RegisterSubscriptionRoutes(router *gin.Engine, db *gorm.DB) {
	// route to create a new task
	PushGroup := router.Group("/push")
	PushGroup.Use(middlewares.CheckCurrentUser())
	PushGroup.POST("/subscribe", func(c *gin.Context) {
		handlers.SubscribePush(c, db)
	})
	PushGroup.POST("/test", func(c *gin.Context) {
		handlers.TestPushHandler(c, db)
	})

}
