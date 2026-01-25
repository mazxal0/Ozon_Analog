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

func (r *AuthRepository) DB() *gorm.DB {
	return r.db
}

// ===================== RefreshToken =====================

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

// ===================== User =====================

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

func (r *AuthRepository) GetUserByIDTx(tx *gorm.DB, ID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := tx.Where("id = ?", ID).First(&user).Error; err != nil {
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

func (r *AuthRepository) UpdateUserCartID(userID, cartID uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("cart_id", cartID).Error
}

// ===================== EmailConfirmation =====================

// Создание токена
func (r *AuthRepository) CreateEmailToken(token *models.EmailConfirmation) error {
	return r.db.Create(token).Error
}

// Получаем валидный код (без Tx)
func (r *AuthRepository) GetValidEmailCode(code, email, codeType string) (*models.EmailConfirmation, error) {
	var token models.EmailConfirmation
	err := r.db.Where(`
		code = ? AND email = ? AND type = ? AND used = false AND expires_at > ?
	`, code, email, codeType, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// Получаем валидный код внутри транзакции
func (r *AuthRepository) GetValidEmailCodeTx(tx *gorm.DB, code, email, codeType string) (*models.EmailConfirmation, error) {
	var token models.EmailConfirmation
	err := tx.Where(`
		code = ? AND email = ? AND type = ? AND used = false AND expires_at > ?
	`, code, email, codeType, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// Ставим EmailVerified = true
func (r *AuthRepository) VerifyEmail(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"email_verified": true,
		"updated_at":     time.Now(),
	}).Error
}

// Помечаем код как использованный внутри транзакции
func (r *AuthRepository) MarkCodeUsedTx(tx *gorm.DB, codeID uuid.UUID) error {
	return tx.Model(&models.EmailConfirmation{}).
		Where("id = ? AND used = false", codeID).
		Updates(map[string]interface{}{
			"used":       true,
			"updated_at": time.Now(),
		}).Error
}

// Инвалидируем все старые коды данного типа внутри транзакции
func (r *AuthRepository) InvalidateCodesTx(tx *gorm.DB, email, codeType string) error {
	return tx.Model(&models.EmailConfirmation{}).
		Where("email = ? AND type = ? AND used = false", email, codeType).
		Updates(map[string]interface{}{
			"used":       true,
			"updated_at": time.Now(),
		}).Error
}

// Счётчик для rate-limit
func (r *AuthRepository) CountCodesByEmail(email string, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.EmailConfirmation{}).
		Where("email = ? AND created_at >= ?", email, since).
		Count(&count).Error
	return count, err
}

func (r *AuthRepository) GetLatestValidEmailCodeTx(tx *gorm.DB, email, codeType string) (*models.EmailConfirmation, error) {
	var token models.EmailConfirmation
	err := tx.
		Where("email = ? AND type = ? AND used = false AND expires_at > ?", email, codeType, time.Now()).
		Order("created_at DESC").
		Limit(1).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}
