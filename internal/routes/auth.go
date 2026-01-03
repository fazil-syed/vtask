package routes

// add routes for handling tasks such as creating a task, marking a task as done, querying tasks
import (
	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/handlers"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(router *gin.Engine, db *gorm.DB) {
	// route to create a new task
	router.POST("/auth/register", func(c *gin.Context) {
		handlers.RegisterUserHandler(c, db)
	})
	// route to get all tasks
	router.POST("/auth/login", func(c *gin.Context) {
		handlers.LoginUserHandler(c, db)
	})
	// router.GET("/auth/logout", func(ctx *gin.Context) {

	// })
	// router.GET("/auth/me", func(ctx *gin.Context) {

	// })
}
