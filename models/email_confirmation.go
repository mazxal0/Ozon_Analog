package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailConfirmation struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID
	User   User

	Email     string `gorm:"not null;index"`         // 👈 для rate-limit
	Type      string `gorm:"size:16;not null;index"` // login | register
	Code      string `gorm:"size:6;not null"`
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

func (e *EmailConfirmation) BeforeSave(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
