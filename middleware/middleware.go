package middleware

import (
	"context"
	"os"
	"strings"
	"time"

	"todolist/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func VerifyToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access denied. No token provided.",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var blacklisted bson.M
	err := config.DB.Collection("blacklist_tokens").FindOne(ctx, bson.M{"token": tokenString}).Decode(&blacklisted)
	if err == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token sudah tidak valid (logout).",
		})
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Invalid or expired token.",
		})
	}

	if id, ok := claims["id"].(string); ok {
		c.Locals("userID", id)
	}
	if role, ok := claims["role"].(string); ok {
		c.Locals("role", role)
	}

	return c.Next()
}

func RequireRole(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Akses ditolak. Role tidak ditemukan.",
			})
		}

		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Akses ditolak. Role tidak sesuai.",
		})
	}
}
