package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/middlewares"
	"github.com/syed.fazil/vtask/internal/routes"
	"gorm.io/gorm"
)

type Config struct {
	DB            *gorm.DB
	UploadPath    string
	MaxUploadSize int64
}

func New(cfg Config) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	router.Use(middlewares.MaxUploadSizeMiddleware(int64(cfg.MaxUploadSize)))
	// register task routes
	routes.RegisterTaskRoutes(router, cfg.DB)
	// register auth routes
	routes.RegisterAuthRoutes(router, cfg.DB)
	routes.RegisterSubscriptionRoutes(router, cfg.DB)

	return router
}
