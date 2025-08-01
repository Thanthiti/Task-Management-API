package auth

import (
	"mymodule/pkg/auth"
	"mymodule/pkg/logger"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	
)
func Middleware(jwtManager auth.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("jwt")
		if tokenStr == "" {
			logger.Log.Warn("Missing JWT cookie")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		token, err := jwtManager.VerifyToken(tokenStr)
		if err != nil || !token.Valid {
			logger.Log.Warn("Invalid JWT token: ", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Log.Error("Failed to parse JWT claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		userID := uint(claims["userID"].(float64))
		logger.Log.Info("Authorized user ID from cookie: ", userID)

		c.Locals("userID", userID)
		return c.Next()
	}
}


