package dto

type AllProcessorsResponseDTO struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	RetailPrice    float64 `json:"retail_price"`
	WholesalePrice float64 `json:"wholesale_price"`
	ImageURL       string  `json:"image_url,omitempty"`
}
