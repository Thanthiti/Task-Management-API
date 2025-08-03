package handler

import (
	// "mymodule/internal/task/model"
	"fmt"
	"mymodule/internal/task/model"
	"mymodule/internal/task/usecase"
	"mymodule/pkg/auth"
	"mymodule/pkg/logger"
	"mymodule/pkg/middleware"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpTaskhandler struct {
	usecase usecase.TaskUsecase
	token  auth.TokenService
	valid   *validator.Validate
}

func NewTaskHandler(app *fiber.App, usecase usecase.TaskUsecase,token auth.TokenService, valid *validator.Validate) {
	handler := &HttpTaskhandler{
		usecase: usecase,
		token: token,
		valid:   valid,
	}
	task := app.Group("/task",middleware.Middleware(token))
	task.Post("/", handler.Create)

}
func (h *HttpTaskhandler) Create(c *fiber.Ctx) error {
	var input model.CreateTaskRequest
	if err := c.BodyParser(&input); err != nil {
		logger.Log.Error("Invalid task request : ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.valid.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// Check panic from c.locals
	rawUserID := c.Locals("userID")
	userIDStr := fmt.Sprintf("%v", rawUserID) 
	userID64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uint(userID64)

	task := model.ToTask(input, uint(userID))
	if err := h.usecase.Create(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Create task successfully "})
}
func (h *HttpTaskhandler) GetTaskByID(c *fiber.Ctx) error {

	return nil
}
func (h *HttpTaskhandler) GetTaskByUser(c *fiber.Ctx) error {
	return nil
}
func (h *HttpTaskhandler) GetTaskByIDAndUser(c *fiber.Ctx) error {
	return nil
}

func (h *HttpTaskhandler) UpdateTask(c *fiber.Ctx) error {
	return nil

}
func (h *HttpTaskhandler) DeleteTask(c *fiber.Ctx) error {
	return nil
}
