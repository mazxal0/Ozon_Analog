package dto

import (
	"eduVix_backend/internal/common/types"
	"github.com/google/uuid"
)

type CreateOrderItem struct {
	ProductID   uuid.UUID         `gorm:"type:uuid;primaryKey"`
	ProductType types.ProductType `gorm:"type:product_type;not null"`
	Quantity    int
	UnitPrice   float64
}
