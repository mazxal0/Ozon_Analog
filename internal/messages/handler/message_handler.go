package handler

import (
	"Market_backend/internal/messages/service"
	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	MessageService *service.MessageService
}

func NewMessageHandler(msgService *service.MessageService) *MessageHandler {
	return &MessageHandler{MessageService: msgService}
}

// Отправка сообщения
func (h *MessageHandler) SendMessage(c *fiber.Ctx) error {
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Text  string `json:"text"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.MessageService.SendMessage(body.Name, body.Email, body.Phone, body.Text); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Message sent successfully"})
}

// Получение всех сообщений
func (h *MessageHandler) GetMessages(c *fiber.Ctx) error {
	messages, err := h.MessageService.GetMessages()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"messages": messages})
}
