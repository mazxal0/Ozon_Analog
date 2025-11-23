package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TwoFactorCode struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null"`        // ссылка на пользователя
	User      *User          `gorm:"foreignKey:UserID"`         // связь с пользователем
	Code      string         `gorm:"type:varchar(10);not null"` // сам код (например, 6-значный)
	ExpiresAt time.Time      `gorm:"not null"`                  // время, когда код перестаёт быть валидным
	Used      bool           `gorm:"default:false"`             // был ли код использован
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // мягкое удаление (опционально)
}
