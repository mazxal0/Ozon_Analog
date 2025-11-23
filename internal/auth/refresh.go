package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"eduVix_backend/models"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	Token     string    `json:"token"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

func GenerateRefreshToken() (*RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	raw := hex.EncodeToString(b)

	hash := sha256.Sum256([]byte(raw))

	return &RefreshToken{
		Token:     raw,
		TokenHash: hex.EncodeToString(hash[:]),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}, nil
}

func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func ValidateRefreshToken(db *gorm.DB, raw string) (*models.RefreshToken, error) {
	hash := HashRefreshToken(raw)
	var rt models.RefreshToken
	if err := db.Where("token_hash = ?", hash).First(&rt).Error; err != nil {
		return nil, err
	}
	if rt.ExpiresAt.Before(time.Now()) {
		db.Delete(&rt)
		return nil, errors.New("refresh token expired")
	}
	return &rt, nil
}
