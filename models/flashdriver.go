package models

import (
	"time"

	"github.com/google/uuid"
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
	USBInterface    string
	FormFactor      string
	ReadSpeed       int
	WriteSpeed      int
	ChipType        string
	OTGSupport      bool
	BodyMaterial    string
	Color           string
	WaterResistance bool
	DustResistance  bool
	Shockproof      bool
	CapType         string // <-- ВАЖНО, добавляем!

	LengthMM    float64
	WidthMM     float64
	ThicknessMM float64
	WeightG     float64

	Compatibility   string
	OperatingTemp   string
	StorageTemp     string
	CountryOfOrigin string
	PackageContents string
	WarrantyMonths  int
	Features        string

	Images []Image `gorm:"foreignKey:FlashDriveID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
