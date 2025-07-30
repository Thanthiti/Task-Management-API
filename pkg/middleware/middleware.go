package auth

import (
	"mymodule/pkg/auth"
	"mymodule/pkg/logger"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	
)

func Middleware(jwtManager auth.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Log.Warn("Missing Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			logger.Log.Warn("Invalid Authorization format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		tokenStr := parts[1]
		token, err := jwtManager.VerifyToken(tokenStr)
		if err != nil || !token.Valid {
			logger.Log.Warn("Invalid token: ", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Log.Error("Cannot parse JWT claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		userID := uint(claims["userID"].(float64))
		logger.Log.Info("Authorized user ID: ", userID)

		c.Locals("userID", userID)

		return c.Next()
	}
}
