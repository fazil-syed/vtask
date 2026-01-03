package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("ksdghksjghkjghkjasdgh")

type JWTClaim struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(email *string, username *string, userID uint) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := JWTClaim{
		Email:    *email,
		UserName: *username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Audience:  jwt.ClaimStrings{"vtask"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtSecret)
	return
}
func ValidateToken(signedToken string) (tokenData *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return
	}
	tokenData, ok := token.Claims.(*JWTClaim)
	if ok && token.Valid {
		return
	}
	err = errors.New("invalid token")
	return
}
