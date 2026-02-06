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
	taskGroup.POST("", func(c *gin.Context) {
		handlers.CreateTaskHandler(c, db)
	})
	// route to get all tasks
	taskGroup.GET("", func(c *gin.Context) {
		handlers.GetTasksHandler(c, db)
	})
	taskGroup.PATCH("/mark-complete/:task_id", func(ctx *gin.Context) {
		handlers.MarkTaskCompletedHandler(ctx, db)
	})
	taskGroup.PATCH("/mark-incomplete/:task_id", func(ctx *gin.Context) {
		handlers.MarkTaskInCompletedHandler(ctx, db)
	})
	taskGroup.DELETE("/delete/:task_id", func(ctx *gin.Context) {
		handlers.DeleteTaskHandler(ctx, db)
	})
	taskGroup.PATCH("/edit/:task_id", func(ctx *gin.Context) {
		handlers.EditTaskHandler(ctx, db)
	})
}
