package models

import (
	"github.com/google/uuid"
	"time"
)

type FlashDrive struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	SKU             string
	Name            string
	Brand           string
	RetailPrice     float64
	WholesalePrice  float64
	WholesaleMinQty int
	Stock           int
	CapacityGB      int
	USBType         string
	ReadSpeedMB     int
	WriteSpeedMB    int
	Features        string
	CountryOfOrigin string
	Images          []Image `gorm:"foreignKey:FlashDriveID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
