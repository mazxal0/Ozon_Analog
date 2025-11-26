package models

import (
	"eduVix_backend/internal/common/types"
	"github.com/google/uuid"
	"time"
)

type OrderItem struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID uuid.UUID
	Order   Order

	ProductID   uuid.UUID         `gorm:"type:uuid;not null"`
	ProductType types.ProductType `gorm:"type:product_type;not null"`

	Quantity  int
	UnitPrice float64 // цена за единицу на момент заказа (опт или розница)
	CreatedAt time.Time
	UpdatedAt time.Time
}
