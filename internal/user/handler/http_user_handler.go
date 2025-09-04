package handler

import (
	"mymodule/internal/user/model"
	"mymodule/internal/user/usecase"
	"mymodule/pkg/auth"
	"mymodule/pkg/logger"
	"mymodule/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpUserhandler struct {
	usecase usecase.UserUsecase
	token   auth.TokenService
	valid   *validator.Validate
}

func NewUserHandler(app *fiber.App, usecase usecase.UserUsecase, token auth.TokenService, valid *validator.Validate) {
	handler := &HttpUserhandler{
		usecase: usecase,
		token:   token,
		valid:   valid,
	}

	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	user := app.Group("/user", middleware.Middleware(token))
	user.Get("/profile", handler.Profile)
	user.Put("/", handler.Updateuser)
	user.Delete("/", handler.DeleteUser)
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

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   3600, // 1 Hours
	})

	return c.JSON(fiber.Map{"token": token})
}

func (h *HttpUserhandler) Profile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	user, err := h.usecase.Profile(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	res := model.ToUserProfileResponse(user)
	return c.JSON(res)
}

func (h *HttpUserhandler) Updateuser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	input := new(model.UpdateUserRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.valid.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := model.User{
		ID:    userID,
		Name:  input.Name,
		Email: input.Email,
	}

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

func (h *HttpUserhandler) Logout(c *fiber.Ctx) error {
	// Clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   -1, // Delete cookie
	})
	return c.JSON(fiber.Map{"message": "logout success"})
}
