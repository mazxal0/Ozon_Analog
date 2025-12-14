package router

import (
	"Market_backend/internal/middleware"
	"Market_backend/internal/payment/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterPaymentRouter(app *fiber.App, h *handler.PaymentHandler) {
	payments := app.Group("/payments")

	// 2. Вебхук ЮKassa
	payments.Post("/webhook", h.Webhook)

	// 1. Создание платежа для фронтенда
	payments.Post("/:order_id", middleware.AuthRequired(), h.CreatePayment)

}
