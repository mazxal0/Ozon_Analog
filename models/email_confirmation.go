package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (e *EmailConfirmation) BeforeSave(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
