package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null"`
	Items     []CartItem `gorm:"foreignKey:CartID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
