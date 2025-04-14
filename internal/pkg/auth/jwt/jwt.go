package jwt

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/auth"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptySecret = errors.New("secret cannot be empty")
)

type Service struct {
	secret []byte
}

func NewService(secretKey []byte) (*Service, error) {
	if len(secretKey) == 0 {
		return nil, ErrEmptySecret
	}
	return &Service{secret: secretKey}, nil
}

func (s *Service) Generate(role entity.UserRole) (string, error) {
	claims := auth.Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) Validate(tokenString string) (*auth.Claims, error) {
	claims := &auth.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, auth.ErrInvalidToken
	}
	if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return claims, nil
}
