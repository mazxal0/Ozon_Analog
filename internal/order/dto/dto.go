package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderItemDTO struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"` // UnitPrice
}

type OrderDTO struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	Status      string         `json:"status"`
	OrderNumber int32          `json:"number"`
	Total       float64        `json:"total"`
	Items       []OrderItemDTO `json:"items"`
	Name        string         `json:"name"`
}

type OrderAdminDTO struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	Status      string         `json:"status"`
	OrderNumber int32          `json:"order_number"`
	Total       float64        `json:"total"`
	Items       []OrderItemDTO `json:"items"`
	Name        string         `json:"name"`
	Surname     string         `json:"surname"`
	LastName    string         `json:"last_name"`
	Number      string         `json:"number"`
	Email       string         `json:"email"`
	LenItems    int            `json:"len_items"`
}

type AllOrdersResponse struct {
	TotalOrders int        `json:"total_orders"`
	TotalItems  int        `json:"total_items"`
	TotalSum    float64    `json:"total_sum"`
	Orders      []OrderDTO `json:"orders"`
}

type AllOrdersAdminResponse struct {
	TotalOrders int             `json:"total_orders"`
	TotalItems  int             `json:"total_items"`
	TotalSum    float64         `json:"total_sum"`
	Orders      []OrderAdminDTO `json:"orders"`
}
