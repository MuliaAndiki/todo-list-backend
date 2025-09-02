package routes

import (
	"todolist/controllers"
	"todolist/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")
	ctrl := controllers.AuthController{}

	auth.Post("/register", ctrl.Register)
	auth.Post("/login", ctrl.Login)
	auth.Post("/logout", ctrl.Logout)
	auth.Get("/getProfile", middleware.VerifyToken, ctrl.GetProfile)
}
