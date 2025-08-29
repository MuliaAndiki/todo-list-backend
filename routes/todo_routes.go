package routes

import (
	"boilerpad/controllers"
	"boilerpad/middleware"

	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(api fiber.Router){
	todo := api.Group("/todo")
	
	todo.Post("/create", middleware.VerifyToken, controllers.TodosController{}.CreateTodos)
	todo.Put("/edit/:_id", controllers.TodosController{}.EditTodos)
	todo.Delete("/delete/:_id", controllers.TodosController{}.HapusTodos)

}