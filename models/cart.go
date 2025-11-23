package models

import (
	"github.com/google/uuid"
	"time"
)

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID
	User      User
	Items     []CartItem `gorm:"foreignKey:CartID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
