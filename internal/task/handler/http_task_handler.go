package handler

import (
	// "mymodule/internal/task/model"
	"mymodule/internal/task/model"
	"mymodule/internal/task/usecase"
	"mymodule/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpTaskhandler struct {
	usecase usecase.TaskUsecase
	valid   *validator.Validate
}

func NewTaskHandler(app *fiber.App, usecase usecase.TaskUsecase, valid *validator.Validate) {
	handler := &HttpTaskhandler{
		usecase: usecase,
		valid:   valid,
	}
	task := app.Group("/task")
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
	userID, ok := rawUserID.(uint)
	if !ok {
		logger.Log.Error("Invalid userID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

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
