package handler

import (
	"eduVix_backend/internal/auth/dto"
	"eduVix_backend/internal/auth/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (handler *AuthHandler) Login(c *fiber.Ctx) error {
	var body dto.AuthLogin
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	access, refresh, err := handler.service.Login(body)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   false, // включи если HTTPS
		SameSite: "Strict",
		Path:     "/auth/",
		MaxAge:   30 * 24 * 60 * 60,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": access,
		"message":      "user login successfully",
	})
}

func (handler *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.AuthRegister

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	access, refresh, err := handler.service.Registration(req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   false, // включи если HTTPS
		SameSite: "Strict",
		Path:     "/auth/",
		MaxAge:   30 * 24 * 60 * 60,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": access,
		"message":      "user registered successfully",
	})
}

func (handler *AuthHandler) Refresh(c *fiber.Ctx) error {
	// Получаем refresh token из cookie
	refreshToken := c.Cookies("refresh_token") // имя cookie
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token not provided",
		})
	}

	// Вызываем сервис для обновления токенов
	accessToken, newRefreshToken, err := handler.service.RefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Обновляем cookie с новым refresh token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HTTPOnly: true,
		Secure:   false, // включить для HTTPS
		SameSite: "Strict",
		Path:     "/auth/",
		MaxAge:   30 * 24 * 60 * 60,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
		"message":      "token refreshed successfully",
	})
}

func (handler *AuthHandler) Logout(c *fiber.Ctx) error {
	refresh := c.Cookies("refresh_token")
	if refresh == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token missing",
		})
	}

	err := handler.service.Logout(refresh)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Удаляем cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/auth/",
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).SendString("Missing token")
	}

	if err := h.service.ConfirmEmail(token); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.SendString("Email verified successfully!")
}
