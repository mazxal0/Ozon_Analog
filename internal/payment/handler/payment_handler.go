package handler

import (
	"Market_backend/internal/common/types"
	"Market_backend/internal/common/utils"
	"Market_backend/internal/order/repository"
	"Market_backend/internal/payment/service"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"os"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	orderRepo      *repository.OrderRepository
}

func NewPaymentHandler(paymentService *service.PaymentService, orderRepo *repository.OrderRepository) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService, orderRepo: orderRepo}
}

// Эндпоинт для фронтенда
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	orderIDStr := c.Params("order_id")

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid order ID",
		})
	}

	userID, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	method := c.Query("method")

	order, err := h.orderRepo.GetOrderById(orderID, userID) // можно добавить userId из JWT
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	var paymentMethod types.PaymentMethod
	switch method {
	case "card":
		paymentMethod = types.PaymentMethodCard
	case "sbp":
		paymentMethod = types.PaymentMethodSBP
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Unknown payment method"})
	}
	payment, confirmationURL, err := h.paymentService.CreatePayment(order, paymentMethod)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"payment_id":       payment.ID,
		"confirmation_url": confirmationURL,
	})
}

type YooKassaWebhook struct {
	Event  string `json:"event"`
	Object struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"object"`
}

func (h *PaymentHandler) Webhook(c *fiber.Ctx) error {
	secret := os.Getenv("YKASSA_WEBHOOK_SECRET")
	body := c.Body()
	signature := c.Get("X-Request-Signature-SHA256")

	// Проверка подписи
	if !verifyYooKassaSignature(secret, body, signature) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid signature")
	}

	// Парсинг JSON
	var webhook YooKassaWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	paymentID := webhook.Object.ID
	status := types.PaymentStatus(webhook.Object.Status) // "succeeded", "canceled", "pending"

	// Обновление статуса платежа через сервис
	if err := h.paymentService.UpdatePaymentStatus(paymentID, status); err != nil {
		fmt.Printf("Failed to update payment status: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update payment")
	}

	return c.SendStatus(fiber.StatusOK)
}

// verifyYooKassaSignature проверяет подпись webhook
func verifyYooKassaSignature(secret string, body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
