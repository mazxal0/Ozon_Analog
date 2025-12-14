package models

import (
	"time"

	"Market_backend/internal/common/types"
	"github.com/google/uuid"
)

type Payment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID
	Order     Order
	Method    types.PaymentMethod `gorm:"type:payment_method;not null"`        // "bank_card", "sbp"
	Status    types.PaymentStatus `gorm:"type:payment_status;default:pending"` // pending, succeeded, canceled
	Amount    float64
	Currency  string
	PaymentID string // ID платежа в ЮKassa
	CreatedAt time.Time
	UpdatedAt time.Time
}
