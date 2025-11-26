package handler

import (
	"eduVix_backend/internal/cart/dto"
	"eduVix_backend/internal/cart/service"
	"eduVix_backend/internal/common/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) AddNewItem(c *fiber.Ctx) error {
	var cartItem dto.CartItemDto
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	cartId, err := h.service.AddNewItem(cartItem)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Item with ID %v was added successfully", cartId),
	})
}

func (h *CartHandler) GetAllCartItems(c *fiber.Ctx) error {
	userId, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err,
		})
	}

	cartIdStr := c.Params("cart_id")
	if cartIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cart_id is required",
		})
	}

	cartId, err := uuid.Parse(cartIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	items, err := h.service.GetAllCartItems(userId, cartId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items": items,
	})

}

func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userId, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err,
		})
	}

	itemIdStr := c.Params("item_id")
	if itemIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "item_id is required",
		})
	}
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	err = h.service.RemoveItem(itemId, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Item with id = %v removed successfully", itemId),
	})
}

func (h *CartHandler) ChangeQuantity(c *fiber.Ctx) error {
	userId, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err,
		})
	}

	quantityStr := c.Query("quantity")
	if quantityStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "quantity is required",
		})
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	itemIdStr := c.Params("item_id")
	if itemIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "item_id is required",
		})
	}
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	err = h.service.ChangeItem(itemId, userId, quantity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Item with ID = %v changed quantity successfully", itemId),
	})
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userId, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err,
		})
	}

	cartIdStr := c.Query("cart_id")
	if cartIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cart_id is required",
		})
	}
	cartId, err := uuid.Parse(cartIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	err = h.service.ClearCart(userId, cartId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cleared cart",
	})
}
