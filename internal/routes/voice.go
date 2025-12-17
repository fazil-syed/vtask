package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/handlers"
	"gorm.io/gorm"
)

func RegisterVoiceRoutes(router *gin.Engine, db *gorm.DB, sttProvider handlers.STTService, uploadPath string) {
	// route to upload a voice note
	router.POST("/voice", func(c *gin.Context) {
		handlers.UploadVoiceNoteHandler(c, db, sttProvider, uploadPath)
	})
}
