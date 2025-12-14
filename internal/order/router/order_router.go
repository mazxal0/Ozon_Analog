package router

import (
	"Market_backend/internal/middleware"
	"Market_backend/internal/order/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterOrderRouter(app *fiber.App, h *handler.OrderHandler) {
	order := app.Group("/order")

	order.Post("/create", middleware.AuthRequired(), h.CreateOrder)
	order.Get("/", middleware.AuthRequired(), h.GetOrderById)

	order.Get("/get-all-orders", middleware.AuthRequired(), h.GetOrders)
	order.Get("/get-for-all-users", middleware.AuthRequired(), middleware.AdminOnly(), h.GetAllOrders)

	order.Post("/update-status", middleware.AuthRequired(), middleware.AdminOnly(), h.UpdateOrderStatusHandler)

}
