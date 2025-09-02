package routes

import (
	"todolist/controllers"
	"todolist/middleware"

	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(api fiber.Router) {
	todo := api.Group("/todo")
	todo.Use(middleware.VerifyToken)

	todo.Post("/create", controllers.TodosController{}.CreateTodos)
	todo.Put("/edit/:_id", controllers.TodosController{}.EditTodos)
	todo.Delete("/delete/:_id", controllers.TodosController{}.HapusTodos)
	todo.Get("/getAll", controllers.TodosController{}.GetAllTodo)
	todo.Get("/getByUser/:userId", controllers.TodosController{}.GetTodoByUser)
}
