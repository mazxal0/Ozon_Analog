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
	Code      string    `gorm:"size:6;not null"` // ✅ 6-значный код
	ExpiresAt time.Time // ✅ время жизни кода
	Used      bool      // ✅ использован или нет
	CreatedAt time.Time
}

func (e *EmailConfirmation) BeforeSave(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
