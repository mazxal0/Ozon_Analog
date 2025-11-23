package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {

		role := c.Locals("role")

		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		if role != "admin" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "forbidden: admin only",
			})
		}

		return c.Next()
	}
}
