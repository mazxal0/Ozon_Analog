package models

import (
	"Market_backend/internal/common/types"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderNumber int32     `gorm:"autoincrement;not null"`
	UserID      uuid.UUID
	User        User
	Status      types.OrderStatus `gorm:"type:order_status;default:in_progress"` // in_progress, paid, delivered
	Total       float64           // общая сумма заказа
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
}
