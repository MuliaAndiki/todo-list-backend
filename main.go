package main

import (
	"log"
	"os"

	"todolist/config"
	"todolist/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found")
	}
	app := fiber.New()
	config.ConnectDB()
	api := app.Group("/go")
	routes.AuthRoutes(api)
	routes.TodoRoutes(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen(":" + port))
}
