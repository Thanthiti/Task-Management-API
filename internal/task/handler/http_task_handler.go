package handler

import (
	"mymodule/internal/task/model"
	"mymodule/internal/task/usecase"
	"mymodule/pkg/auth"
	"mymodule/pkg/helper"
	"mymodule/pkg/logger"
	"mymodule/pkg/middleware"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpTaskhandler struct {
	usecase usecase.TaskUsecase
	token   auth.TokenService
	valid   *validator.Validate
}

func NewTaskHandler(app *fiber.App, usecase usecase.TaskUsecase, token auth.TokenService, valid *validator.Validate) {
	handler := &HttpTaskhandler{
		usecase: usecase,
		token:   token,
		valid:   valid,
	}
	task := app.Group("/task", middleware.Middleware(token))
	task.Post("/", handler.Create)
	task.Get("/", handler.GetTaskByUser)
	task.Get("/:id", handler.GetTaskByIDAndUser)
	task.Put("/:id", handler.UpdateTask)
	task.Delete("/:id", handler.DeleteTask)

	// for Admin get all task regardless userID
	task.Get("/admin/:id", handler.GetTaskByID)

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
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		logger.Log.Error("Unauthorized access: ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	task := model.ToTask(input, uint(userID))
	if err := h.usecase.Create(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Create task successfully "})
}

// For admin
func (h *HttpTaskhandler) GetTaskByID(c *fiber.Ctx) error {

	return nil
}

// All task
func (h *HttpTaskhandler) GetTaskByUser(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		logger.Log.Error("Unauthorized access: ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	tasks, err := h.usecase.GetByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch tasks"})
	}
	var resp []model.TaskResponse
	if tasks != nil {
		resp = model.ToTaskResponseList(*tasks)
	}
	return c.JSON(resp)

}

// Detail task
func (h *HttpTaskhandler) GetTaskByIDAndUser(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		logger.Log.Error("Unauthorized access: ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	taskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid task ID"})
	}

	task, err := h.usecase.GetByIDAndUser(uint(taskID), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found or unauthorized"})
	}
	var resp model.TaskResponse
	if task != nil {
		resp = model.ToTaskResponse(*task)
	}

	return c.JSON(resp)
}
	
func (h *HttpTaskhandler) UpdateTask(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		logger.Log.Error("Unauthorized access: ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	taskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid task ID"})
	}

	var input model.UpdateTaskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.valid.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.usecase.UpdateTask(&input,uint(taskID), uint(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "task updated"})
}


func (h *HttpTaskhandler) DeleteTask(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		logger.Log.Error("Unauthorized access: ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	taskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid task ID"})
	}
	if err := h.usecase.DeleteTask(uint(taskID),userID); err != nil{
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "task deleted"})
}
