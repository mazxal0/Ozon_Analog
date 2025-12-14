package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	Phone     string
	Text      string `gorm:"not null"`
	CreatedAt time.Time
}
