package models

import (
	"Market_backend/internal/common/validate"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"Market_backend/internal/common/types"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string
	Surname      string
	LastName     string
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	Number       string
	Role         types.UserRole `gorm:"type:user_role;default:'user'"`
	AvatarURL    string

	CartID uuid.UUID `gorm:"type:uuid;not null;unique"`

	RefreshToken  *RefreshToken  `gorm:"constraint:OnDelete:CASCADE;"`
	TwoFactorCode *TwoFactorCode `gorm:"constraint:OnDelete:CASCADE;"`

	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if err := u.Validate(); err != nil {
		return err
	}

	return nil
}

//func (u *User) BeforeUpdate(tx *gorm.DB) error {
//	if err := u.Validate(); err != nil {
//		return err
//	}
//	return nil
//}

func (u *User) Validate() error {
	return validate.Validate.Struct(u)
}
