package router

import (
	"eduVix_backend/internal/cart/handler"
	"eduVix_backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterCartRouter(app *fiber.App, h *handler.CartHandler) {
	cart := app.Group("/cart")

	cart.Post("/", middleware.AuthRequired(), h.AddNewItem)
	cart.Get("/:cart_id", middleware.AuthRequired(), h.GetAllCartItems)
}
