package handler

import (
	"Market_backend/internal/common/utils"
	"Market_backend/internal/user/dto"
	"Market_backend/internal/user/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAllUsers GET /users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	return h.GetAllUsers(c)
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.service.GetMe(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) ChangeMe(c *fiber.Ctx) error {
	var userChange dto.UserChange

	if err := c.BodyParser(&userChange); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.service.ChangeMe(userChange, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
