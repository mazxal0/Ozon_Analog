package router

import (
	"eduVix_backend/internal/middleware"
	"eduVix_backend/internal/product/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterFlashDriverRouter(app *fiber.App, h *handler.FlashDriveHandler) {
	fd := app.Group("/flash-driver")

	fd.Post("/", middleware.AuthRequired(), h.CreateFlashDrive)
	fd.Delete("/:flashId", middleware.AuthRequired(), h.DeleteFlashDrive)

	fd.Get("/", middleware.AuthRequired(), h.GetAllFlashDrives)
	fd.Get("/:flashId", middleware.AuthRequired(), h.GetFlashDriveById)
	fd.Patch("/:flashId", middleware.AuthRequired(), h.UpdateFlashDrive)
}
