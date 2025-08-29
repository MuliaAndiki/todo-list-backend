package controllers

import (
	"boilerpad/config"
	"boilerpad/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodosController struct{}

func(TodosController) CreateTodos(c *fiber.Ctx) error{
	var body struct{
		Todos string `json:"todos"`
	}
	if err := c.BodyParser(&body); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Invalid Req Body",
		})
	}

	todo := models.Todo{
		Todos: body.Todos,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := config.DB.Collection("todos")
	res, err := collection.InsertOne(c.Context(), todo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create todo",
		})
	}

	todo.ID = res.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message" : "Todo Berhasil Dibuat",
		"data": todo,
	})
}

func(TodosController) EditTodos(c *fiber.Ctx) error{
	idParams := c.Params("_id")
	objID, err := primitive.ObjectIDFromHex(idParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID format",
		})
	}

	var body struct{
		Todos string `json:"todos"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	collection := config.DB.Collection("todos")
	update := bson.M{
		"$set": bson.M{
			"todos" : body.Todos,
			"update_at" : time.Now(),
		},
	}

	result, err := collection.UpdateOne(c.Context(),bson.M{"_id":objID},update)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON((fiber.Map{
			"message" : "Invalid IdParams",
		}))
	} 

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message" : "Todo Not Found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo updated successfully",
		"data" : update,
	})
}

func(TodosController) HapusTodos(c *fiber.Ctx) error{
	idParams := c.Params("_id")
	objId, err := primitive.ObjectIDFromHex(idParams)

	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "ID Params Invalid",
		})
	}
	collection := config.DB.Collection("todos")
	res, err := collection.DeleteOne(c.Context(),bson.M{"_id" : objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete todo",
		})
	}
	if res.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message" : "Todo Not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo berhasil dihapus",
	})
}