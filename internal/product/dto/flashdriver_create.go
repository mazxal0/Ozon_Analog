package dto

import "mime/multipart"

// DTO для создания
type FlashDriveCreateDTO struct {
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Brand string `json:"brand"`

	RetailPrice     float64 `json:"retail_price"`
	WholesalePrice  float64 `json:"wholesale_price"`
	WholesaleMinQty int     `json:"wholesale_min_qty"`
	Stock           int     `json:"stock"`

	CapacityGB      int    `json:"capacity_gb"`
	USBInterface    string `json:"usb_interface"` // USB 2.0 / USB 3.0 / USB 3.2 etc.
	FormFactor      string `json:"form_factor"`
	ReadSpeed       int    `json:"read_speed"`
	WriteSpeed      int    `json:"write_speed"`
	ChipType        string `json:"chip_type"`
	OTGSupport      bool   `json:"otg_support"`
	BodyMaterial    string `json:"body_material"`
	Color           string `json:"color"`
	WaterResistance bool   `json:"water_resistance"`
	DustResistance  bool   `json:"dust_resistance"`
	Shockproof      bool   `json:"shockproof"`
	CapType         string `json:"cap_type"`

	LengthMM    float64 `json:"length_mm"`
	WidthMM     float64 `json:"width_mm"`
	ThicknessMM float64 `json:"thickness_mm"`
	WeightG     float64 `json:"weight_g"`

	Compatibility   string `json:"compatibility"`
	OperatingTemp   string `json:"operating_temp"`
	StorageTemp     string `json:"storage_temp"`
	CountryOfOrigin string `json:"country_of_origin"`
	PackageContents string `json:"package_contents"`
	WarrantyMonths  int    `json:"warranty_months"`
	Features        string `json:"features"`

	ImageFiles []*multipart.FileHeader `json:"image_files"`
}
