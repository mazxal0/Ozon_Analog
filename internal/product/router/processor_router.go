package router

import (
	"Market_backend/internal/middleware"
	"Market_backend/internal/product/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterProcessorRouter(app *fiber.App, h *handler.ProcessorHandler) {
	proc := app.Group("/processor")

	proc.Post("/", middleware.AuthRequired(), h.CreateProcessor)
	proc.Delete("/:procId", middleware.AuthRequired(), h.DeleteProcessor)

	proc.Get("/", h.GetAllProcessors)
	proc.Get("/:procId", h.GetProcessorById)
	proc.Patch("/:procId", middleware.AuthRequired(), h.UpdateProcessor)
}
