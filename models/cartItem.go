package models

import (
	"Market_backend/internal/common/types"
	"github.com/google/uuid"
	"time"
)

type CartItem struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	CartID uuid.UUID
	Cart   Cart

	ProductID   uuid.UUID         `gorm:"type:uuid;not null"`
	ProductType types.ProductType `gorm:"type:product_type;not null"`

	//ProcessorID  uuid.UUID
	//Processor    *Processor
	//FlashDriveID uuid.UUID
	//FlashDrive   *FlashDrive
	Quantity  int
	UnitPrice float64 // текущая цена на момент добавления в корзину
	CreatedAt time.Time
	UpdatedAt time.Time
}
