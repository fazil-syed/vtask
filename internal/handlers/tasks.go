package handlers

import (
	"log"
	"net/http"
	"strconv"
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
	strUserId, exists := c.Get("user_id")

	if !exists {
		log.Println("User id not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	userId, ok := strUserId.(uint)
	if !ok {
		log.Println("Invalid user")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	task := models.Task{
		Name:      taskSchema.Title,
		Completed: false,
		Content:   taskSchema.Content,
		DueAt:     dueAt,
		UserId:    userId,
	}
	if err := db.WithContext(c.Request.Context()).Create(&task).Error; err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// create reminders
	for _, r := range taskSchema.Reminders {
		var remindAt time.Time
		if r.OffsetMinutes != nil && task.DueAt != nil {
			remindAt = task.DueAt.Add(-time.Duration(*r.OffsetMinutes) * time.Minute)
		} else {
			continue
		}
		if err := db.WithContext(c.Request.Context()).Create(&models.Reminder{TaskID: task.ID, RemindAt: remindAt}).Error; err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}

// Handler to get all tasks
func GetTasksHandler(c *gin.Context, db *gorm.DB) {
	strUserId, exists := c.Get("user_id")

	if !exists {
		log.Println("User id not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	userId, ok := strUserId.(uint)
	if !ok {
		log.Println("Invalid user")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	var tasks []models.Task
	if err := db.WithContext(c.Request.Context()).Where("user_id = ?", userId).Order("id DESC").Find(&tasks).Error; err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	response := make([]schemas.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, taskToResponse(task))
	}
	c.JSON(http.StatusOK, response)
}

func taskToResponse(task models.Task) schemas.TaskResponse {
	var dueAt *string
	if task.DueAt != nil {
		s := task.DueAt.Format(time.RFC3339)
		dueAt = &s
	}

	return schemas.TaskResponse{
		ID:        task.ID,
		Title:     task.Name,
		Content:   task.Content,
		DueAt:     dueAt,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	}
}

func MarkTaskCompletedHandler(c *gin.Context, db *gorm.DB) {
	strUserId, exists := c.Get("user_id")

	if !exists {
		log.Println("User id not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	userId, ok := strUserId.(uint)
	if !ok {
		log.Println("Invalid user")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	taskId := c.Param("task_id")
	intTaskId, err := strconv.Atoi(taskId)
	if err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	var task models.Task
	if err := db.WithContext(c.Request.Context()).Where("user_id = ? AND id = ?", userId, intTaskId).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	isCompleted := task.Completed
	if isCompleted {
		c.JSON(http.StatusOK, gin.H{"message": "Task already marked as completed"})
		return
	}
	err = db.WithContext(c.Request.Context()).Model(&task).Update("completed", true).Error
	if err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as completed"})
	return

}

func MarkTaskInCompletedHandler(c *gin.Context, db *gorm.DB) {
	strUserId, exists := c.Get("user_id")

	if !exists {
		log.Println("User id not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	userId, ok := strUserId.(uint)
	if !ok {
		log.Println("Invalid user")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	taskId := c.Param("task_id")
	intTaskId, err := strconv.Atoi(taskId)
	if err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	var task models.Task
	if err := db.WithContext(c.Request.Context()).Where("user_id = ? AND id = ?", userId, intTaskId).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	isCompleted := task.Completed
	if !isCompleted {
		c.JSON(http.StatusOK, gin.H{"message": "Task already marked as not completed"})
		return
	}
	err = db.WithContext(c.Request.Context()).Model(&task).Update("completed", false).Error
	if err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as not completed"})
	return

}

func DeleteTaskHandler(c *gin.Context, db *gorm.DB) {
	strUserId, exists := c.Get("user_id")

	if !exists {
		log.Println("User id not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	userId, ok := strUserId.(uint)
	if !ok {
		log.Println("Invalid user")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	taskId := c.Param("task_id")
	intTaskId, err := strconv.Atoi(taskId)
	if err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}
	res := db.WithContext(c.Request.Context()).Where("user_id = ? AND id = ?", userId, intTaskId).Delete(&models.Task{})
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
	return

}

func EditTaskHandler(c *gin.Context, db *gorm.DB) {
	strUserId, exists := c.Get("user_id")

	if !exists {
		log.Println("User id not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	userId, ok := strUserId.(uint)
	if !ok {
		log.Println("Invalid user")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	taskId := c.Param("task_id")
	intTaskId, err := strconv.Atoi(taskId)
	if err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}
	var taskSchema schemas.EditTaskInput
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

	var task models.Task
	if err := db.WithContext(c.Request.Context()).Where("user_id = ? AND id = ?", userId, intTaskId).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	updates := map[string]interface{}{}
	if dueAt != nil && dueAt != task.DueAt {
		updates["due_at"] = dueAt
	}
	if taskSchema.Content != "" && taskSchema.Content != task.Content {
		updates["content"] = taskSchema.Content
	}
	if taskSchema.Title != "" && taskSchema.Title != task.Name {
		updates["name"] = taskSchema.Title
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No changes to update"})
		return
	}
	err = db.WithContext(c.Request.Context()).Model(&task).Updates(updates).Error
	if err != nil {
		log.Printf("ERROR : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
	return

}
