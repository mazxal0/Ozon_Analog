package handler

import (
	"Market_backend/internal/common/types"
	"Market_backend/internal/common/utils"
	"Market_backend/internal/order/dto"
	"Market_backend/internal/order/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userId, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
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
			"error": err.Error(),
		})
	}

	orderId, err := h.service.CreateOrder(userId, cartId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"order_id": orderId,
	})
}

func (h *OrderHandler) GetOrderById(c *fiber.Ctx) error {
	userId, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	orderIdStr := c.Query("order_id")
	if orderIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order_id is required",
		})
	}
	orderId, err := uuid.Parse(orderIdStr)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	order, err := h.service.GetOrderById(userId, orderId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"order": order,
	})
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {

	userId, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	orders, err := h.service.GetAllUserOrders(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := dto.AllOrdersResponse{
		TotalOrders: len(orders),
	}

	var ordersDTO []dto.OrderDTO

	for _, order := range orders {
		var itemsDTO []dto.OrderItemDTO

		for _, item := range order.Items {

			var name string

			// Получаем название из разных таблиц
			switch item.ProductType {
			case types.Processor:
				proc, err := h.service.ProcService.GetProcessorById(item.ProductID)
				if err != nil {
					name = "Unknown processor"
				} else {
					name = proc.Name
				}

			case types.FlashDriver:
				fd, err := h.service.FlashService.GetFlashDriveById(item.ProductID)
				if err != nil {
					name = "Unknown flash drive"
				} else {
					name = fd.Name
				}

			default:
				name = "Unknown product"
			}

			totalPrice := item.UnitPrice * float64(item.Quantity)

			response.TotalItems += item.Quantity
			if order.Status != types.Cancelled {
				response.TotalSum += totalPrice
			}

			itemsDTO = append(itemsDTO, dto.OrderItemDTO{
				Name:     name,
				Quantity: item.Quantity,
				Price:    item.UnitPrice,
			})
		}

		ordersDTO = append(ordersDTO, dto.OrderDTO{
			ID:        order.ID,
			CreatedAt: order.CreatedAt,
			Status:    string(order.Status),
			Total:     order.Total,
			Items:     itemsDTO,
		})
	}

	response.Orders = ordersDTO

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.service.GetAllOrders() // используем существующий метод сервиса
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := dto.AllOrdersResponse{}

	var ordersDTO []dto.OrderDTO

	for _, order := range orders {
		var itemsDTO []dto.OrderItemDTO
		var orderItemsCount int
		var totalSum float64

		for _, item := range order.Items {
			var name string
			switch item.ProductType {
			case types.Processor:
				proc, err := h.service.ProcService.GetProcessorById(item.ProductID)
				if err != nil {
					name = "Unknown processor"
				} else {
					name = proc.Name
				}
			case types.FlashDriver:
				fd, err := h.service.FlashService.GetFlashDriveById(item.ProductID)
				if err != nil {
					name = "Unknown flash drive"
				} else {
					name = fd.Name
				}
			default:
				name = "Unknown product"
			}

			totalPrice := item.UnitPrice * float64(item.Quantity)
			orderItemsCount += item.Quantity
			totalSum += totalPrice

			itemsDTO = append(itemsDTO, dto.OrderItemDTO{
				Name:     name,
				Quantity: item.Quantity,
				Price:    item.UnitPrice,
			})
		}

		ordersDTO = append(ordersDTO, dto.OrderDTO{
			ID:        order.ID,
			Name:      order.User.Name, // нужно чтобы User был preloaded
			Status:    string(order.Status),
			Total:     totalSum,
			CreatedAt: order.CreatedAt,
			Items:     itemsDTO,
		})

		response.TotalItems += orderItemsCount
		response.TotalSum += totalSum
	}

	response.TotalOrders = len(ordersDTO)
	response.Orders = ordersDTO

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *OrderHandler) UpdateOrderStatusHandler(c *fiber.Ctx) error {
	// Получаем order_id и новый статус из query
	orderIDStr := c.Query("order_id")
	newStatusStr := c.Query("status")

	if orderIDStr == "" || newStatusStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order_id and status are required",
		})
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order_id",
		})
	}

	var newStatus types.OrderStatus
	switch newStatusStr {
	case "in_progress":
		newStatus = types.InProgress
	case "paid":
		newStatus = types.Paid
	case "completed":
		newStatus = types.Completed
	case "cancelled":
		newStatus = types.Cancelled
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid status",
		})
	}

	// Меняем статус
	if err := h.service.UpdateOrderStatus(orderID, newStatus); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "status updated",
	})
}
