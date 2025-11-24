package models

import (
	"eduVix_backend/internal/common/validate"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"eduVix_backend/internal/common/types"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name         string         `validate:"required,max=50,min=2"`
	Surname      string         `validate:"required,max=50,min=2"`
	LastName     string         `validate:"required,max=50,min=2"`
	Email        string         `gorm:"uniqueIndex" validate:"required,email"`
	PasswordHash string         `validate:"required,password"`
	Number       string         `validate:"omitempty,number"`
	Role         types.UserRole `gorm:"type:user_role;default:'user'"`
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	RefreshToken  *RefreshToken  `gorm:"constraint:OnDelete:CASCADE;"`
	TwoFactorCode *TwoFactorCode `gorm:"constraint:OnDelete:CASCADE;"`

	EmailVerified bool
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if err := u.Validate(); err != nil {
		fmt.Println("Validation error:", err)
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
