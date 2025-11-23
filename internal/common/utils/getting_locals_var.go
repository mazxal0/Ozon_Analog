package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUserId(c *fiber.Ctx) (uuid.UUID, error) {
	raw := c.Locals("userId")
	if raw == nil {
		return uuid.Nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found in token",
		})
	}
	userIDStr, ok := raw.(string)
	if !ok {
		return uuid.Nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID has invalid type in token",
		})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID format",
		})
	}
	
	return userID, nil
}
