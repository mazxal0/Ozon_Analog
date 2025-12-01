package repository

import (
	"Market_backend/internal/common"
	"Market_backend/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		db: common.DB,
	}
}

func (r *AuthRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) GetRefreshTokenByHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	if err := r.db.Where("token_hash = ?", hash).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) UpdateRefreshToken(token *models.RefreshToken) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(&token).Error
}

func (r *AuthRepository) DeleteRefreshTokenByHash(hash string) error {
	return r.db.Where("token_hash = ?", hash).Delete(&models.RefreshToken{}).Error
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByID(ID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) DeleteUserByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.User{}).Error
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) CreateEmailToken(token *models.EmailConfirmation) error {
	return r.db.Create(token).Error
}

// 1️⃣ Получаем токен подтверждения, который действителен
func (r *AuthRepository) GetValidEmailToken(tokenStr string) (*models.EmailConfirmation, error) {
	var token models.EmailConfirmation
	err := r.db.Preload("User").
		Where("token = ? AND used = false AND expires_at > ?", tokenStr, time.Now()).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// 2️⃣ Ставим EmailVerified = true для пользователя
func (r *AuthRepository) VerifyEmail(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("email_verified", true).Error
}

// 3️⃣ Помечаем токен как использованный
func (r *AuthRepository) MarkTokenUsed(tokenID uuid.UUID) error {
	return r.db.Model(&models.EmailConfirmation{}).Where("id = ?", tokenID).Update("used", true).Error
}
