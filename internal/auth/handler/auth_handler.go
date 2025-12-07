package handler

import (
	"Market_backend/internal/auth/dto"
	"Market_backend/internal/auth/service"

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

	// Токены не возвращаем, только проверка пароля + отправка кода на email
	if err := handler.service.Login(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login code sent to your email",
	})
}

func (handler *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.AuthRegister

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	// Токены не возвращаем, только регистрация + отправка кода
	if err := handler.service.Registration(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "confirmation code sent to your email",
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

func (handler *AuthHandler) ConfirmCode(c *fiber.Ctx) error {
	type ConfirmDTO struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	var dto ConfirmDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	access, refresh, err := handler.service.ConfirmCode(dto.Code, dto.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ставим cookie для refresh token
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
		"message":      "user authenticated successfully",
	})
}
