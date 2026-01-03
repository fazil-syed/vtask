package routes

// add routes for handling tasks such as creating a task, marking a task as done, querying tasks
import (
	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/handlers"
	"github.com/syed.fazil/vtask/internal/middlewares"
	"gorm.io/gorm"
)

func RegisterTaskRoutes(router *gin.Engine, db *gorm.DB) {
	// route to create a new task

	taskGroup := router.Group("/tasks")
	taskGroup.Use(middlewares.CheckCurrentUser())
	taskGroup.POST("/", func(c *gin.Context) {
		handlers.CreateTaskHandler(c, db)
	})
	// route to get all tasks
	taskGroup.GET("/", func(c *gin.Context) {
		handlers.GetTasksHandler(c, db)
	})
}
