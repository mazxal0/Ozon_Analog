package dto

import (
	"github.com/google/uuid"
	"time"
)

type OrderItemDTO struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"` // UnitPrice
}

type OrderDTO struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	Status    string         `json:"status"`
	Total     float64        `json:"total"`
	Items     []OrderItemDTO `json:"items"`
	Name      string         `json:"name"`
}

type AllOrdersResponse struct {
	TotalOrders int        `json:"total_orders"`
	TotalItems  int        `json:"total_items"`
	TotalSum    float64    `json:"total_sum"`
	Orders      []OrderDTO `json:"orders"`
}
