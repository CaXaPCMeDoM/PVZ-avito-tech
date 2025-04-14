package dummy

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/auth"
)

type AuthUseCase struct {
	jwtService auth.TokenService
}

func NewDummyAuthUseCase(jwtService auth.TokenService) *AuthUseCase {
	return &AuthUseCase{jwtService: jwtService}
}

func (d *AuthUseCase) GenerateDummyToken(role entity.UserRole) (string, error) {
	return d.jwtService.Generate(role)
}
