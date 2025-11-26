package router

import (
	"eduVix_backend/internal/middleware"
	"eduVix_backend/internal/order/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterOrderRouter(app *fiber.App, h *handler.OrderHandler) {
	order := app.Group("/order")

	order.Post("/create", middleware.AuthRequired(), h.CreateOrder)
}
