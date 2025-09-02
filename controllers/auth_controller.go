package controllers

import (
	"context"
	"os"
	"strings"
	"time"

	"todolist/config"
	"todolist/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (AuthController) Register(c *fiber.Ctx) error {
	var body struct {
		Fullname string       `json:"fullname"`
		Email    string       `json:"email"`
		Password string       `json:"password"`
		Role     *models.Role `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	var existing models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": body.Email}).Decode(&existing)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email sudah terdaftar",
		})
	}

	var role models.Role

	if body.Role == nil {
		role = models.RoleUser
	} else {
		role = *body.Role
		if !role.IsValid() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid role, must be one of: admin, user",
			})
		}
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	user := models.User{
		ID:        primitive.NewObjectID(),
		Fullname:  body.Fullname,
		Email:     body.Email,
		Password:  string(hash),
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = config.DB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Register berhasil",
		"user": fiber.Map{
			"id":       user.ID.Hex(),
			"fullname": user.Fullname,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (AuthController) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": body.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Email tidak ditemukan",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Password salah",
		})
	}

	claims := jwt.MapClaims{
		"id":    user.ID.Hex(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, _ := token.SignedString([]byte(secret))

	return c.JSON(fiber.Map{

		"message": "Login berhasil",
		"data": fiber.Map{
			"token": tokenString,
			"user": fiber.Map{
				"id":       user.ID.Hex(),
				"fullname": user.Fullname,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}

func (AuthController) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak ditemukan",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.DB.Collection("blacklist_tokens").InsertOne(ctx, bson.M{
		"token":     tokenString,
		"createdAt": time.Now(),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal logout",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logout berhasil, token di-blacklist",
	})
}

func (AuthController) GetProfile(c *fiber.Ctx) error {
	userId, ok := c.Locals("userID").(string)
	if !ok || userId == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: no user claims found",
		})
	}

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	collection := config.DB.Collection("users")
	var user models.User
	err = collection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get Profile Success",
		"data": fiber.Map{
			"id":        user.ID.Hex(),
			"fullname":  user.Fullname,
			"email":     user.Email,
			"role":      user.Role,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}
