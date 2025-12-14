package router

import (
	"Market_backend/internal/messages/handler"
	"Market_backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterMessageRoutes(app *fiber.App, h *handler.MessageHandler) {
	message := app.Group("/message")

	// Отправка сообщения (POST /message/send)
	message.Post("/send", h.SendMessage)

	// Получение всех сообщений (GET /message/all) с middleware если нужно
	message.Get("/all", middleware.AuthRequired(), middleware.AdminOnly(), h.GetMessages)
}
