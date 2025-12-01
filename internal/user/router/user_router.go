package router

import (
	"Market_backend/internal/middleware"
	"Market_backend/internal/user/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App, h *handler.UserHandler) {
	user := app.Group("/users")

	// GET /users
	user.Get("/", middleware.AuthRequired(), middleware.AdminOnly(), h.GetAllUsers)

	// GET /users/me
	user.Get("/me", middleware.AuthRequired(), h.GetMe)
}
