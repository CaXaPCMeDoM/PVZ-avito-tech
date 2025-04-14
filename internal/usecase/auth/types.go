package auth

import (
	"PVZ-avito-tech/internal/entity"
	"github.com/google/uuid"
)

type LoginResponse struct {
	Role entity.UserRole `json:"role"`
}

type RegisterResponse struct {
	Id    uuid.UUID       `json:"id"`
	Email string          `json:"email"`
	Role  entity.UserRole `json:"role"`
}
