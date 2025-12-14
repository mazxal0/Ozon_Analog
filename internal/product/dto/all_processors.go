package dto

import "github.com/google/uuid"

type AllProcessorsResponseDTO struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	RetailPrice    float64   `json:"retail_price"`
	WholesalePrice float64   `json:"wholesale_price"`
	ImageURL       *string   `json:"image_url,omitempty" gorm:"column:image_url"`
}
