package models

import (
	"github.com/google/uuid"
	"time"
)

type EmailConfirmation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID
	User      User
	Token     string    `gorm:"uniqueIndex"` // уникальный токен
	ExpiresAt time.Time // время действия токена
	Used      bool      // был ли использован
	CreatedAt time.Time
}
