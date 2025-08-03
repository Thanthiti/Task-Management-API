package helper

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetUserIDFromContext(c *fiber.Ctx) (uint, error) {
	rawUserID := c.Locals("userID")
	if rawUserID == nil {
		return 0, fmt.Errorf("userID not found in context")
	}

	userIDStr := fmt.Sprintf("%v", rawUserID)
	userID64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid userID format: %w", err)
	}

	return uint(userID64), nil
}