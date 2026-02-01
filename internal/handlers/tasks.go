package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/schemas"
	"gorm.io/gorm"
)

// Handler to create a new task
func CreateTaskHandler(c *gin.Context, db *gorm.DB) {
	// use schema from schemas package to bind the request body
	// and use the models package for the Task model
	var taskSchema schemas.CreateTaskInput
	if err := c.ShouldBindJSON(&taskSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var dueAt *time.Time
	if taskSchema.DueDate != nil {

		loc, err := time.LoadLocation(taskSchema.Timezone)
		if err != nil {
			log.Printf("ERROR : %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timezone"})
			return
		}
		t, err := time.ParseInLocation("2006-01-02 15:04", *taskSchema.DueDate, loc)
		if err != nil {
			log.Printf("ERROR : %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		dueAt = &t
	}
	task := models.Task{
		Name:      taskSchema.Title,
		Completed: false,
		Content:   taskSchema.Content,
		DueAt:     dueAt,
	}
	if err := db.WithContext(c.Request.Context()).Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// Handler to get all tasks
func GetTasksHandler(c *gin.Context, db *gorm.DB) {
	var tasks []models.Task
	if err := db.WithContext(c.Request.Context()).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
