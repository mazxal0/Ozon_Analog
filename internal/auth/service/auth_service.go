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

func (s *AuthService) Registration(dto dto.AuthRegister) error {
	if err := dto.Validate(); err != nil {
		return err
	}
	if dto.Password != dto.RepeatPassword {
		return errors.New("passwords do not match")
	}

	user, err := s.repo.GetUserByEmail(dto.Email)
	if err != nil {
		return err
	}
	if user != nil {
		return errors.New("user already exists")
	}

	hashPassword, _ := utils.HashPassword(dto.Password)

	return s.repo.DB().Transaction(func(tx *gorm.DB) error { // <-- здесь напрямую Transaction
		userID := uuid.New()
		cartID := uuid.New()

		user := &models.User{
			ID:            userID,
			Email:         dto.Email,
			Name:          dto.Name,
			Surname:       dto.Surname,
			LastName:      dto.LastName,
			Number:        dto.Number,
			Role:          dto.Role,
			PasswordHash:  hashPassword,
			EmailVerified: false,
			CartID:        cartID, // сразу присваиваем CartID
		}

		cart := &models.Cart{
			ID:     cartID,
			UserID: userID,
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if err := tx.Create(cart).Error; err != nil {
			return err
		}

		// Создание EmailConfirmation
		code := utils.GenerateSixDigitCode()
		confirmation := &models.EmailConfirmation{
			UserID:    userID,
			Code:      code,
			ExpiresAt: time.Now().Add(10 * time.Minute),
			Used:      false,
		}

		if err := tx.Create(confirmation).Error; err != nil {
			return err
		}

		// Отправка письма
		body := fmt.Sprintf(`
<h1>Market</h1>
<p>Здравствуйте, %s!</p>
<p>Ваш код подтверждения:</p>
<h2>%s</h2>
<p>Код действует 10 минут.</p>
`, user.Name, code)

		return s.mailSender.SendEmail(user.Email, "Подтвердите почту", body)
	})
}

func (s *AuthService) Login(dto dto.AuthLogin) error {
	if err := dto.Validate(); err != nil {
		return err
	}

	user, err := s.repo.GetUserByEmail(dto.Email)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user does not exist")
	}

	// ✅ Проверка пароля
	if !utils.CheckPasswordHash(dto.Password, user.PasswordHash) {
		return errors.New("invalid user data")
	}

	// ✅ Генерируем код для ВХОДА
	loginCode := utils.GenerateSixDigitCode()

	confirmation := &models.EmailConfirmation{
		UserID:    user.ID,
		Code:      loginCode,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Used:      false,
	}

	if err := s.repo.CreateEmailToken(confirmation); err != nil {
		return err
	}

	// ✅ Отправляем код на почту
	body := fmt.Sprintf(`
<h1>Market</h1>
<p>Здравствуйте, %s!</p>
<p>Код для входа:</p>
<h2>%s</h2>
<p>Код действует 10 минут.</p>
`,
		user.Name,
		loginCode,
	)

	if err := s.mailSender.SendEmail(
		user.Email,
		"Код для входа",
		body,
	); err != nil {
		return err
	}

	// ✅ ТУТ НЕТ ТОКЕНОВ!
	return nil
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

	accessToken, err := auth.GenerateToken(user.ID.String(), string(user.Role), user.Name, user.CartID.String())
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken.Token, nil
}

func (s *AuthService) ConfirmCode(code, email string) (string, string, error) {
	confirmation, err := s.repo.GetValidEmailCode(code, email)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUserByID(confirmation.UserID)
	if err != nil {
		return "", "", err
	}

	// Если email не подтверждён — это регистрация
	if !user.EmailVerified {
		if err := s.repo.VerifyEmail(user.ID); err != nil {
			return "", "", err
		}
	}

	// Помечаем код как использованный
	if err := s.repo.MarkCodeUsed(confirmation.ID); err != nil {
		return "", "", err
	}

	// Генерируем токены
	access, refreshToken, err := s.issueTokens(user)
	if err != nil {
		return "", "", err
	}

	return access, refreshToken, nil
}

func (s *AuthService) issueTokens(user *models.User) (string, string, error) {
	access, err := auth.GenerateToken(user.ID.String(), string(user.Role), user.Name, user.CartID.String())
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
