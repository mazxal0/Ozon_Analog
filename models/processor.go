package models

import (
	"github.com/google/uuid"
	"time"
)

type Processor struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	SKU                string
	Name               string
	Brand              string
	RetailPrice        float64 // цена для розницы
	WholesalePrice     float64 // цена для опта
	WholesaleMinQty    int     // минимальное количество для опта
	Stock              int
	Line               string
	Architecture       string
	Socket             string
	BaseFrequency      float64
	TurboFrequency     float64
	Cores              int
	Threads            int
	L1Cache            string
	L2Cache            string
	L3Cache            string
	Lithography        string
	TDP                int
	Features           string
	MemoryType         string
	MaxRAM             string
	MaxRAMFrequency    string
	IntegratedGraphics bool
	GraphicsModel      string
	MaxTemperature     int
	PackageContents    string
	CountryOfOrigin    string
	Images             []Image `gorm:"foreignKey:ProcessorID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
