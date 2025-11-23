package dto

import (
	"eduVix_backend/internal/common/types"
	"github.com/google/uuid"
)

type CartItemDto struct {
	CartID      uuid.UUID         `json:"cart_id"`
	ProductId   uuid.UUID         `json:"product_id"`
	ProductType types.ProductType `json:"product_type"`
	Quantity    int               `json:"quantity"`
}
