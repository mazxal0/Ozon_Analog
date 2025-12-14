package dto

import (
	"Market_backend/internal/common/types"
	"Market_backend/models"
	"github.com/google/uuid"
)

type GetCartItemsResponse struct {
	ID          uuid.UUID         `json:"id"`
	ProductId   uuid.UUID         `json:"product_id"`
	ProductType types.ProductType `json:"product_type"`
	Quantity    int               `json:"quantity"`
	ImageUrl    string            `json:"image_url"`
	Price       float64           `json:"price"`

	Name string `json:"name"`
}

type CartItemWithProduct struct {
	CartItem models.CartItem
	Product  any
}
