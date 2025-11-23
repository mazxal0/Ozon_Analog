package models

import (
	"github.com/google/uuid"
	"time"
)

type Image struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProcessorID  *uuid.UUID
	FlashDriveID *uuid.UUID
	URL          string
	CreatedAt    time.Time
}
