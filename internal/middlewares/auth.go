package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syed.fazil/vtask/internal/utils"
)

func CheckCurrentUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		signedToken, err := ctx.Cookie("auth_token")
		if err != nil {
			// Cookie is missing
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie missing"})
			return
		}
		tokenData, err := utils.ValidateToken(signedToken)
		if err != nil {
			ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		ctx.Set("user_email", tokenData.Email)
		ctx.Set("user_name", tokenData.UserName)
		ctx.Set("user_id", tokenData.UserID)
		ctx.Next()
	}
}
