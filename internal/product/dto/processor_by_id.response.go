package dto

import (
	"github.com/google/uuid"
)

type ProcessorWithImagesDTO struct {
	ID                 uuid.UUID `json:"id"`
	SKU                string    `json:"sku"`
	Name               string    `json:"name"`
	Brand              string    `json:"brand"`
	RetailPrice        float64   `json:"retail_price"`
	WholesalePrice     float64   `json:"wholesale_price"`
	WholesaleMinQty    int       `json:"wholesale_min_qty"`
	Stock              int       `json:"stock"`
	Line               string    `json:"line"`
	Architecture       string    `json:"architecture"`
	Socket             string    `json:"socket"`
	BaseFrequency      float64   `json:"base_frequency"`
	TurboFrequency     float64   `json:"turbo_frequency"`
	Cores              int       `json:"cores"`
	Threads            int       `json:"threads"`
	L1Cache            string    `json:"l1_cache"`
	L2Cache            string    `json:"l2_cache"`
	L3Cache            string    `json:"l3_cache"`
	Lithography        string    `json:"lithography"`
	TDP                int       `json:"tdp"`
	Features           string    `json:"features"`
	MemoryType         string    `json:"memory_type"`
	MaxRAM             string    `json:"max_ram"`
	MaxRAMFrequency    string    `json:"max_ram_frequency"`
	IntegratedGraphics bool      `json:"integrated_graphics"`
	GraphicsModel      string    `json:"graphics_model"`
	MaxTemperature     int       `json:"max_temperature"`
	PackageContents    string    `json:"package_contents"`
	CountryOfOrigin    string    `json:"country_of_origin"`
	CountOrders        int       `json:"count_orders"`
	ImageURLs          []string  `json:"image_urls"` // только URL
}
