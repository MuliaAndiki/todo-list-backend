package utils

import (
	"os"
	"time"

	"todolist/models"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(user models.User) (string, error) {
	// Fallback secret jika JWT_SECRET tidak di-set
	if len(JwtSecret) == 0 {
		JwtSecret = []byte("fallback-secret-key-change-in-production")
	}

	claims := jwt.MapClaims{
		"id":    user.ID.Hex(),
		"email": user.Email,
		"role":  string(user.Role),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
