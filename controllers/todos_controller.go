package controllers

import (
	"time"

	"todolist/config"
	"todolist/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodosController struct{}

func (TodosController) CreateTodos(c *fiber.Ctx) error {
	var body struct {
		Todos     string        `json:"todos"`
		Status    models.Status `json:"status"`
		CreatedAt time.Time     `json:"created_at"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Req Body",
		})
	}

	if body.Status == "" {
		body.Status = models.StatusPending
	}

	if !body.Status.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid status value, must be one of: pending, progress, done",
		})
	}

	if body.CreatedAt.IsZero() {
		body.CreatedAt = time.Now()
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User not authenticated",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid user ID in token",
		})
	}

	todo := models.Todo{
		Todos:     body.Todos,
		Status:    body.Status,
		UserID:    userID,
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
		"message": "Todo Create Succes",
		"data":    todo,
	})
}

func (TodosController) EditTodos(c *fiber.Ctx) error {
	idParams := c.Params("_id")
	objID, err := primitive.ObjectIDFromHex(idParams)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID ",
		})
	}

	var body struct {
		Todos  *string        `json:"todos"`
		Status *models.Status `json:"status"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	updateData := bson.M{}

	if body.Todos != nil {
		updateData["todos"] = *body.Todos
	}

	if body.Status != nil {
		if !(*body.Status).IsValid() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid status value, must be one of: pending, progress, done",
			})
		}
		updateData["status"] = *body.Status
	}

	if len(updateData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nothing to update",
		})
	}

	collection := config.DB.Collection("todos")
	update := bson.M{"$set": updateData}

	result, err := collection.UpdateOne(c.Context(), bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update todo",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Todo Not Found",
		})
	}
	var updateTodo models.Todo
	if err := collection.FindOne(c.Context(), bson.M{"_id": objID}).Decode(&updateTodo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch updated todo",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo updated successfully",
		"data":    updateTodo,
	})
}

func (TodosController) HapusTodos(c *fiber.Ctx) error {
	idParams := c.Params("_id")
	objId, err := primitive.ObjectIDFromHex(idParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID Params Invalid",
		})
	}
	collection := config.DB.Collection("todos")
	res, err := collection.DeleteOne(c.Context(), bson.M{"_id": objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete todo",
		})
	}
	if res.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Todo Not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo berhasil dihapus",
	})
}

func (TodosController) GetAllTodo(c *fiber.Ctx) error {
	collection := config.DB.Collection("todos")

	var todos []models.Todo

	res, err := collection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Server Internal Error",
		})
	}
	defer res.Close(c.Context())

	for res.Next(c.Context()) {
		var todo models.Todo

		if err := res.Decode(&todo); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error Decoding Todo",
			})
		}
		todos = append(todos, todo)
	}

	if err := res.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cursor Error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get Todo Successfully ",
		"data":    todos,
	})
}

func (TodosController) GetTodoByUser(c *fiber.Ctx) error {
	collection := config.DB.Collection("todos")

	userId := c.Params("userId")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format Id Invalid",
		})
	}
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid User ID",
		})
	}
	var todos []models.Todo

	filter := bson.M{"userId": objectID}
	res, err := collection.Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Server Internal Error",
		})
	}
	defer res.Close(c.Context())
	for res.Next(c.Context()) {
		var todo models.Todo
		if err := res.Decode(&todo); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error Decoding Todo",
			})
		}
		todos = append(todos, todo)
	}
	if err := res.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cursor Error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get Todo by User Successfully",
		"data":    todos,
	})
}
