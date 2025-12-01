package service

import (
	"Market_backend/internal/auth"
	"Market_backend/internal/auth/dto"
	"Market_backend/internal/auth/repository"
	"Market_backend/internal/common"
	"Market_backend/internal/common/utils"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	CartRepo "Market_backend/internal/cart/repository"
	"Market_backend/internal/mail/service"

	"Market_backend/models"
	"errors"
	"time"
)

type AuthService struct {
	repo       *repository.AuthRepository
	cartRepo   *CartRepo.CartRepository
	mailSender *mail.MailService
}

func NewAuthService(repo *repository.AuthRepository, cartRepo *CartRepo.CartRepository) *AuthService {
	return &AuthService{repo: repo, cartRepo: cartRepo, mailSender: mail.NewMailService()}
}

func (s *AuthService) Registration(dto dto.AuthRegister) (string, string, error) {
	if err := dto.Validate(); err != nil {
		return "", "", err
	}
	if dto.Password != dto.RepeatPassword {
		return "", "", errors.New("passwords do not match")
	}

	user, err := s.repo.GetUserByEmail(dto.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", err
	}
	if user != nil {
		return "", "", errors.New("user already exists")
	}

	hashPassword, _ := utils.HashPassword(dto.Password)
	newUser := &models.User{
		Email:         dto.Email,
		Name:          dto.Name,
		Surname:       dto.Surname,
		PasswordHash:  hashPassword,
		LastName:      dto.LastName,
		Number:        dto.Number,
		Role:          dto.Role,
		EmailVerified: false,
	}

	if err := s.repo.CreateUser(newUser); err != nil {
		return "", "", err
	}
	if err := s.cartRepo.CreateCart(newUser.ID); err != nil {
		return "", "", err
	}

	// Генерация токена подтверждения email
	emailToken := uuid.New().String()
	confirmation := &models.EmailConfirmation{
		UserID:    newUser.ID,
		Token:     emailToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.CreateEmailToken(confirmation); err != nil {
		return "", "", err
	}

	// Отправка письма
	link := fmt.Sprintf("https://yourdomain.com/verify-email?token=%s", emailToken)
	body := fmt.Sprintf(
		`
		<h1>Gay SHOP</h1>
		<p>Здравствуйте, %s!</p>
		<p>Перейдите по ссылке, чтобы подтвердить email и то что вы Гей))))))) АМЕРИКА USA ()()()() BOOBS:</p>
		<a href="%s">Подтвердить email</a>`,
		newUser.Name, link,
	)

	if err := s.mailSender.SendEmail(newUser.Email, "Подтвердите почту", body); err != nil {
		return "", "", err
	}

	// Генерация JWT и refresh token
	access, err := auth.GenerateToken(newUser.ID.String(), string(newUser.Role))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = s.repo.CreateRefreshToken(&models.RefreshToken{
		UserID:    newUser.ID,
		TokenHash: refreshToken.TokenHash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		return "", "", err
	}

	return access, refreshToken.Token, nil
}

func (s *AuthService) Login(dto dto.AuthLogin) (string, string, error) {
	if err := dto.Validate(); err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUserByEmail(dto.Email)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", errors.New("user does not exist")
	}

	isCheck := utils.CheckPasswordHash(dto.Password, user.PasswordHash)

	if !isCheck {
		return "", "", errors.New("invalid user data")
	}

	access, err := auth.GenerateToken(user.ID.String(), string(user.Role))

	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.GenerateRefreshToken()

	if err != nil {
		return "", "", err
	}

	err = s.repo.UpdateRefreshToken(&models.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshToken.TokenHash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})

	if err != nil {
		return "", "", err
	}

	return access, refreshToken.Token, nil
}

func (s *AuthService) Logout(token string) error {
	hash, err := utils.HashPassword(token)
	if err != nil {
		return err
	}
	return s.repo.DeleteRefreshTokenByHash(hash)
}

func (s *AuthService) RefreshToken(token string) (string, string, error) {
	refreshToken, err := auth.ValidateRefreshToken(common.DB, token)

	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := auth.GenerateRefreshToken()

	if err != nil {
		return "", "", err
	}

	err = s.repo.UpdateRefreshToken(&models.RefreshToken{
		UserID:    refreshToken.UserID,
		TokenHash: newRefreshToken.TokenHash,
		ExpiresAt: newRefreshToken.ExpiresAt,
	})

	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUserByID(refreshToken.UserID)
	if err != nil {
		return "", "", err
	}

	accessToken, err := auth.GenerateToken(user.ID.String(), string(user.Role))
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken.Token, nil
}

func (s *AuthService) ConfirmEmail(token string) error {
	confirmation, err := s.repo.GetValidEmailToken(token)
	if err != nil {
		return err
	}

	if err := s.repo.VerifyEmail(confirmation.UserID); err != nil {
		return err
	}

	return s.repo.MarkTokenUsed(confirmation.ID)
}
