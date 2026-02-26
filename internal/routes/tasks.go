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
	newTaskStream := handlers.GetTaskStream()
	taskHandlers := handlers.NewTaskHandler(newTaskStream, db)
	taskGroup := router.Group("/tasks")
	taskGroup.Use(middlewares.CheckCurrentUser())
	taskGroup.POST("", func(c *gin.Context) {
		taskHandlers.CreateTaskHandler(c)
	})

	// route to get all tasks
	taskGroup.GET("", func(c *gin.Context) {
		taskHandlers.GetTasksHandler(c)
	})
	taskGroup.PATCH("/mark-complete/:task_id", func(ctx *gin.Context) {
		taskHandlers.MarkTaskCompletedHandler(ctx)
	})
	taskGroup.PATCH("/mark-incomplete/:task_id", func(ctx *gin.Context) {
		taskHandlers.MarkTaskInCompletedHandler(ctx)
	})
	taskGroup.DELETE("/delete/:task_id", func(ctx *gin.Context) {
		taskHandlers.DeleteTaskHandler(ctx)
	})
	taskGroup.PATCH("/edit/:task_id", func(ctx *gin.Context) {
		taskHandlers.EditTaskHandler(ctx)
	})
	taskGroup.GET("/stream", newTaskStream.SSEConnMiddleware(), func(ctx *gin.Context) {
		newTaskStream.TaskStreamHandler(ctx)
	})
}
