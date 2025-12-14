package dto

import "github.com/google/uuid"

// DTO для списка (каталога)
type AllFlashDrivesResponseDTO struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	WholesalePrice float64   `json:"wholesale_price"`
	RetailPrice    float64   `json:"retail_price"`
	ImageURL       string    `json:"image_url"`
}
 