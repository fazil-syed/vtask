package server

import (
	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/handlers"
	"github.com/syed.fazil/vtask/internal/middlewares"
	"github.com/syed.fazil/vtask/internal/routes"
	"gorm.io/gorm"
)

func SetupServer(db *gorm.DB, sttProvider handlers.STTService, uploadPath string, maxUploadSize int64) *gin.Engine {
	router := gin.Default()
	// register task routes
	routes.RegisterTaskRoutes(router, db)
	// register voice routes
	routes.RegisterVoiceRoutes(router, db, sttProvider, uploadPath)

	router.Use(middlewares.MaxUploadSizeMiddleware(int64(maxUploadSize)))
	return router
}
