package dto

import (
	"eduVix_backend/internal/common/validate"
)

type AuthLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

func (a *AuthLogin) Validate() error {
	return validate.Validate.Struct(a)
}
