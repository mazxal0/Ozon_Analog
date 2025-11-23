package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID
	User      User
	Status    string  // in_progress, paid, delivered
	Total     float64 // общая сумма заказа
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []OrderItem `gorm:"foreignKey:OrderID"`
}
