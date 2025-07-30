package handler

import (
	"fmt"
	"mymodule/internal/user/model"
	"mymodule/internal/user/usecase"
	auths "mymodule/pkg/auth"
	"mymodule/pkg/logger"
	"mymodule/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpUserhandler struct {
	usecase usecase.UserUsecase
	token   auths.TokenService
	valid   *validator.Validate
}

func NewUserHandler(app *fiber.App, usecase usecase.UserUsecase, token auths.TokenService, valid *validator.Validate) {
	handler := &HttpUserhandler{
		usecase: usecase,
		token:   token,
		valid:   valid,
	}

	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	fmt.Println(token)
	user := app.Group("/user", auth.Middleware(token))
	user.Put("/:id", handler.Updateuser)
	user.Delete("/:id", handler.DeleteUser)
}

func (h *HttpUserhandler) Register(c *fiber.Ctx) error {
	var input model.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		logger.Log.Error("Invalid register request : ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.valid.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	user := model.ToUserModel(input)
	if err := h.usecase.Register(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

func (h *HttpUserhandler) Login(c *fiber.Ctx) error {
	var input model.LoginRequest
	if err := c.BodyParser(&input); err != nil {
		logger.Log.Error("Invalid register request : ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.valid.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	token, err := h.usecase.Login(input.Email, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (h *HttpUserhandler) Updateuser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input model.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.valid.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := model.ToUserModel(input)
	user.ID = userID

	if err := h.usecase.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User updated"})
}

func (h *HttpUserhandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	if err := h.usecase.DeleteUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User deleted"})
}
