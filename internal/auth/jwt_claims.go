package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	CartID string `json:"cart_id"`
	jwt.RegisteredClaims
}
