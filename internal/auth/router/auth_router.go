package router

import (
	"eduVix_backend/internal/auth/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRouter(app *fiber.App, h *handler.AuthHandler) {
	auth := app.Group("/auth")

	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)

	// Новый маршрут для подтверждения email
	// Пользователь кликает по ссылке из письма: /auth/verify-email?token=...
	auth.Get("/verify-email", h.VerifyEmail)
}
