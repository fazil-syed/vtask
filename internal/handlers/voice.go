package handlers

// write handlers/voice.go that uses the STTService interface to transcribe an audio file and return the transcription in the response
import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/parser"
	"github.com/syed.fazil/vtask/internal/utils"
	"gorm.io/gorm"
)

func UploadVoiceNoteHandler(c *gin.Context, db *gorm.DB, sttProvider STTService, uploadPath string) {
	ctx := c.Request.Context()
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	// save the file to the upload path
	if err := c.SaveUploadedFile(file, uploadPath+"/"+file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	manual := strings.TrimSpace(c.PostForm("transcription"))
	var transcription string
	if manual != "" {
		transcription = manual
	} else {
		// transcribe the voice note
		transcription, err = sttProvider.Transcribe(uploadPath + "/" + file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transcribe audio"})
			return
		}
	}
	intent := parser.ParseIntent(ctx, transcription)

	// save voice note metadata to the database
	voiceNote := models.VoiceNote{}
	voiceNote.Filename = file.Filename
	voiceNote.ContentType = file.Header.Get("Content-Type")
	voiceNote.Size = file.Size
	voiceNote.Transcription = transcription
	voiceNote.Status = "transcribed"
	voiceNote.FilePath = uploadPath + "/" + file.Filename
	if err := db.WithContext(c.Request.Context()).Create(&voiceNote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save voice note metadata"})
		return
	}
	action := gin.H{}
	switch intent.Name {
	case "create_task":
		task := models.Task{Name: intent.Title, Completed: false}
		if intent.Reminder != "" {
			if t, err := utils.ParseTimeString(intent.Reminder); err == nil && t != nil {
				task.ReminderAt = t
			}
		}
		if err := db.WithContext(c.Request.Context()).Create(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task from voice note"})
			return
		}
		voiceNote.Status = "task_created"
		if err := db.WithContext(c.Request.Context()).Save(&voiceNote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update voice note status"})
			return
		}
		action = gin.H{"action": "create_task", "task_id": task.ID, "task_name": task.Name}
	case "mark_done":
		var task models.Task
		norm := strings.ToLower(strings.TrimSpace(intent.Title))
		err := db.WithContext(c.Request.Context()).Where("lower(name) LIKE ? and completed = ?", "%"+norm+"%", false).Order("created_at desc").First(&task).Error
		if err != nil {
			voiceNote.Status = "task_not_found"
			if err := db.WithContext(c.Request.Context()).Save(&voiceNote).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update voice note status"})
				return
			}

		} else {
			task.Completed = true
			action = gin.H{"action": "mark_done", "task_id": task.ID, "task_name": task.Name}
			if err := db.WithContext(c.Request.Context()).Save(&task).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark task as done"})
				return
			}
			voiceNote.Status = "task_completed"
			if err := db.WithContext(c.Request.Context()).Save(&voiceNote).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update voice note status"})
				return
			}
		}
	}
	c.JSON(http.StatusCreated, gin.H{
		"voice_note":    voiceNote,
		"transcription": transcription,
		"intent":        intent,
		"action":        action,
	})
}
