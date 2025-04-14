package dto

import "PVZ-avito-tech/internal/entity"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required"`
	Role     entity.UserRole `json:"role" binding:"required"`
}

type DummyLoginRequest struct {
	Role entity.UserRole `json:"role" binding:"required"`
}
