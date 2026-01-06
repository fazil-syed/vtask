package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/config"
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
	// register task routes
	routes.RegisterTaskRoutes(router, cfg.DB)
	// register auth routes
	routes.RegisterAuthRoutes(router, cfg.DB)

	router.Use(middlewares.MaxUploadSizeMiddleware(int64(cfg.MaxUploadSize)))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.App.APPUIBaseURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	return router
}
