package dto

type FlashDriveFilterDTO struct {
	Brands       []string `json:"brands"`
	CapacityGB   []int    `json:"capacity_gb"`
	USBInterface []string `json:"usb_interface"`

	PriceAsc bool `json:"price_asc"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
