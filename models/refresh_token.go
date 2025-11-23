package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	TokenHash string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

//func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {
//	if r.ID == 0 {}
//}
