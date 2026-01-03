package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/schemas"
	"github.com/syed.fazil/vtask/internal/utils"
	"gorm.io/gorm"
)

// Handler to register a new user
func RegisterUserHandler(c *gin.Context, db *gorm.DB) {
	dbWithCtx := db.WithContext(c.Request.Context())
	var userSchema schemas.UserRegisterInputSchema
	if err := c.ShouldBindJSON(&userSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var count int64
	// Check if identiy/user already exists
	err := dbWithCtx.Model(&models.Identity{}).Where("issuer = ? AND email = ?", "password", userSchema.Email).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}
	// Check if username if unique
	if err := dbWithCtx.Model(&models.User{}).Where("user_name = ?", userSchema.UserName).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Username %s is already taken", userSchema.UserName)})
		return
	}

	passwordHash, err := utils.HashPassword(userSchema.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	user := models.User{
		PrimaryEmail: userSchema.Email,
		UserName:     userSchema.UserName,
	}
	if err := dbWithCtx.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userIdentity := models.Identity{
		Issuer:       "password",
		Email:        &userSchema.Email,
		PasswordHash: &passwordHash,
		UserId:       user.ID,
		Subject:      strconv.FormatUint(uint64(user.ID), 10),
	}
	if err := dbWithCtx.Create(&userIdentity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
