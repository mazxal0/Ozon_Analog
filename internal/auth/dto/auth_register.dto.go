package dto

import (
	"eduVix_backend/internal/common/types"
	"eduVix_backend/internal/common/validate"
)

type AuthRegister struct {
	Email          string         `json:"email"`
	Password       string         `json:"password" validate:"required,password"`
	RepeatPassword string         `json:"repeat_password" validate:"required,password"`
	Name           string         `json:"name"`
	Surname        string         `json:"surname"`
	LastName       string         `json:"last_name"`
	Number         string         `json:"number"`
	Role           types.UserRole `gorm:"type:user_role;default:'student'"`
}

func (a *AuthRegister) Validate() error {
	return validate.Validate.Struct(a)
}
