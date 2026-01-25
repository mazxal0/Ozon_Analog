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

// ===================== Registration =====================
func (s *AuthService) Registration(dto dto.AuthRegister) error {
	if err := dto.Validate(); err != nil {
		return err
	}
	if dto.Password != dto.RepeatPassword {
		return errors.New("passwords do not match")
	}

	if user, _ := s.repo.GetUserByEmail(dto.Email); user != nil {
		return errors.New("user already exists")
	}

	hashPassword, _ := utils.HashPassword(dto.Password)

	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
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
			CartID:        cartID,
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Cart{ID: cartID, UserID: userID}).Error; err != nil {
			return err
		}

		if err := s.checkEmailRateLimit(dto.Email); err != nil {
			return err
		}

		// инвалидируем старые коды
		if err := s.repo.InvalidateCodesTx(tx, dto.Email, "register"); err != nil {
			return err
		}

		code := utils.GenerateSixDigitCode()
		confirmation := &models.EmailConfirmation{
			UserID:    userID,
			Email:     dto.Email,
			Type:      "register",
			Code:      code,
			ExpiresAt: time.Now().Add(10 * time.Minute),
			Used:      false,
		}

		if err := tx.Create(confirmation).Error; err != nil {
			return err
		}

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

// ===================== Login =====================
func (s *AuthService) Login(dto dto.AuthLogin) error {
	if err := dto.Validate(); err != nil {
		return err
	}

	user, err := s.repo.GetUserByEmail(dto.Email)
	if err != nil || user == nil {
		return errors.New("invalid user data")
	}

	if !utils.CheckPasswordHash(dto.Password, user.PasswordHash) {
		return errors.New("invalid user data")
	}

	if err := s.checkEmailRateLimit(user.Email); err != nil {
		return err
	}

	loginCode := utils.GenerateSixDigitCode()

	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.repo.InvalidateCodesTx(tx, user.Email, "login"); err != nil {
			return err
		}

		confirmation := &models.EmailConfirmation{
			UserID:    user.ID,
			Email:     user.Email,
			Type:      "login",
			Code:      loginCode,
			ExpiresAt: time.Now().Add(10 * time.Minute),
			Used:      false,
		}

		if err := tx.Create(confirmation).Error; err != nil {
			return err
		}

		body := fmt.Sprintf(`
<h1>Market</h1>
<p>Здравствуйте, %s!</p>
<p>Код для входа:</p>
<h2>%s</h2>
<p>Код действует 10 минут.</p>
`, user.Name, loginCode)

		return s.mailSender.SendEmail(user.Email, "Код для входа", body)
	})
}

// ===================== Logout =====================
func (s *AuthService) Logout(token string) error {
	hash, err := utils.HashPassword(token)
	if err != nil {
		return err
	}
	return s.repo.DeleteRefreshTokenByHash(hash)
}

// ===================== RefreshToken =====================
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

// ===================== ConfirmCode =====================
func (s *AuthService) ConfirmCode(code, email string) (string, string, error) {
	var accessToken, refreshToken string

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		// берём последний валидный код типа login
		confirmation, err := s.repo.GetLatestValidEmailCodeTx(tx, email, "login")
		if err != nil || confirmation.Code != code {
			// если login не подошёл — пробуем register
			confirmation, err = s.repo.GetLatestValidEmailCodeTx(tx, email, "register")
			if err != nil || confirmation.Code != code {
				return errors.New("invalid or expired code")
			}
		}

		user, err := s.repo.GetUserByIDTx(tx, confirmation.UserID)
		if err != nil {
			return err
		}

		// если регистрация — верифицируем email
		if confirmation.Type == "register" && !user.EmailVerified {
			if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Update("email_verified", true).Error; err != nil {
				return err
			}
		}

		// помечаем код как использованный
		if err := s.repo.MarkCodeUsedTx(tx, confirmation.ID); err != nil {
			return err
		}

		// выдаем токены
		accessToken, refreshToken, err = s.issueTokens(user)
		return err
	})

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ===================== issueTokens =====================
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

// ===================== RateLimit =====================
func (s *AuthService) checkEmailRateLimit(email string) error {
	count, err := s.repo.CountCodesByEmail(email, time.Now().Add(-1*time.Hour))
	if err != nil {
		return err
	}

	if count >= 20 {
		return errors.New("слишком много запросов. Попробуйте позже")
	}

	return nil
}
