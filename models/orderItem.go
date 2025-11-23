package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderItem struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID      uuid.UUID
	Order        Order
	ProcessorID  *uuid.UUID
	Processor    *Processor
	FlashDriveID *uuid.UUID
	FlashDrive   *FlashDrive
	Quantity     int
	UnitPrice    float64 // цена за единицу на момент заказа (опт или розница)
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
