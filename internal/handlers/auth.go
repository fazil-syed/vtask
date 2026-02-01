package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/config"
	"github.com/syed.fazil/vtask/internal/models"
	"github.com/syed.fazil/vtask/internal/schemas"
	"github.com/syed.fazil/vtask/internal/services"
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
	_, err = services.FindOrCreateUserWithIdentity(c.Request.Context(),
		db,
		schemas.IdentityInput{
			Issuer:        "password",
			Subject:       userSchema.Email,
			Email:         &userSchema.Email,
			EmailVerified: false,
			PasswordHash:  &passwordHash,
			UserName:      &userSchema.UserName,
		},
		services.IntentRegister,
	)
	if err != nil {
		if errors.Is(err, services.ErrIdentityAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	fmt.Println("handler done")
	c.Status(http.StatusCreated)
	return
}
func LoginUserHandler(c *gin.Context, db *gorm.DB) {
	dbWithCtx := db.WithContext(c.Request.Context())
	var loginUserSchema schemas.UserLoginInputSchema
	if err := c.ShouldBindJSON(&loginUserSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check if user exists
	var identiy models.Identity
	if err := dbWithCtx.Model(&models.Identity{}).Where("issuer = ? AND email = ?", "password", loginUserSchema.Email).Preload("User").First(&identiy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// check if password is correct

	if err := utils.CheckPassword(*identiy.PasswordHash, loginUserSchema.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := utils.GenerateToken(identiy.User.PrimaryEmail, identiy.User.UserName, identiy.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// set the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func LogoutUserHandler(c *gin.Context, db *gorm.DB) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.Status(http.StatusNoContent)
	return
}

func GetUserProfileHandler(c *gin.Context, db *gorm.DB) {
	dbWithCtx := db.WithContext(c.Request.Context())
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	var user models.User
	if err := dbWithCtx.Model(&models.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	profile := schemas.UserProfileDataResponse{
		UserName: user.UserName,
		Email:    user.PrimaryEmail,
	}
	c.JSON(http.StatusOK, profile)
	return
}

func InitiateGoogleSSOAuthHandler(c *gin.Context) {
	oauth2Config, err := utils.NewOauth2Config(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	state := utils.GenerateRandomString()
	c.SetCookie("state", state, 300, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(state))
	return
}

func GoogleSSOCallbackHandler(c *gin.Context, db *gorm.DB) {
	stateFromGoogleQuery := c.Query("state")
	cfg := config.App
	oauth2Config, err := utils.NewOauth2Config(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	stateFromCookie, err := c.Cookie("state")
	if err != nil || stateFromGoogleQuery != stateFromCookie {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid OAuth state",
		})
		return
	}
	oauth2Token, err := oauth2Config.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	provider, err := utils.NewOIDCProvider(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.OAuthGoogleClientID})
	idToken, err := verifier.Verify(c.Request.Context(), rawIDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user, err := services.FindOrCreateUserWithIdentity(
		c.Request.Context(),
		db,
		schemas.IdentityInput{
			Issuer:        "google",
			Subject:       idToken.Subject,
			Email:         &claims.Email,
			EmailVerified: claims.Verified,
		},
		services.IntentLogin,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	// clear state cookie
	c.SetCookie("state", "", -1, "/", "", false, true)
	// Generate auth token
	token, err := utils.GenerateToken(user.PrimaryEmail, user.UserName, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// set the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, cfg.APPUIBaseURL)
	return

}
