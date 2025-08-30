package routes

import (
	"todolist/controllers"
	"todolist/middleware"

	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(api fiber.Router) {
	todo := api.Group("/todo")

	todo.Post("/create", middleware.VerifyToken, controllers.TodosController{}.CreateTodos)
	todo.Put("/edit/:_id", middleware.VerifyToken, controllers.TodosController{}.EditTodos)
	todo.Delete("/delete/:_id", middleware.VerifyToken, controllers.TodosController{}.HapusTodos)
	todo.Get("/getAll", middleware.VerifyToken, controllers.TodosController{}.GetAllTodo)
}
