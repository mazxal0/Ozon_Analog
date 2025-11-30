package dto

type FlashDriveFilterDTO struct {
	Brands     []string `json:"brands"`
	CapacityGB []int    `json:"capacity_gb"`
	USBType    []string `json:"usb_type"`
	USBVersion []string `json:"usb_version"`

	PriceAsc bool `json:"price_asc"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
