package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/schemas"
	"gorm.io/gorm"
)

// Handler to create a new task
func CreateTaskHandler(c *gin.Context, db *gorm.DB) {
	// use schema from schemas package to bind the request body
	// and use the models package for the Task model
	var taskScmea schemas.CreateTaskInput
	if err := c.ShouldBindJSON(&taskScmea); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := models.Task{
		Name:      taskScmea.Name,
		Completed: taskScmea.Completed,
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
