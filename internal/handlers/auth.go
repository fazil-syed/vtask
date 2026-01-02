package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/schemas"
	"github.com/syed.fazil/vtask/internal/utils"
	"gorm.io/gorm"
)

// Handler to create a new task
func RegisterUserHandler(c *gin.Context, db *gorm.DB) {
	var userSchema schemas.UserRegisterInputSchema
	if err := c.ShouldBindJSON(&userSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	passwordHash, err := utils.HashPassword(userSchema.Password)
	if err != nil {
		log.Fatalf("Error while hashing password: %v", err)
	}
	user := models.User{
		Email:        userSchema.Email,
		UserName:     userSchema.UserName,
		PasswordHash: passwordHash,
	}
	if err := db.WithContext(c.Request.Context()).Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
