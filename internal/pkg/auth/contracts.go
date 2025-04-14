package auth

import (
	"PVZ-avito-tech/internal/entity"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Claims struct {
	Role entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type (
	TokenService interface {
		Generate(role entity.UserRole) (string, error)
		Validate(tokenString string) (*Claims, error)
	}
)
