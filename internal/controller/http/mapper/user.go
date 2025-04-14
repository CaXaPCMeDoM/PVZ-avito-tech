package mapper

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
)

func RegisterRequestToEntityUser(body dto.RegisterRequest) *entity.User {
	return &entity.User{
		Email:    body.Email,
		Password: body.Password,
		Role:     body.Role,
	}
}
