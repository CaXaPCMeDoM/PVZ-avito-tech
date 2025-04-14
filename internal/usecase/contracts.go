package usecase

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/usecase/auth"
	"context"
	"github.com/google/uuid"
)

type (
	Auth interface {
		Register(ctx context.Context, u *entity.User) (auth.RegisterResponse, error)
		Login(ctx context.Context, email string, rawPassword string) (auth.LoginResponse, error)
	}
	DummyLogin interface {
		GenerateDummyToken(role entity.UserRole) (string, error)
	}
	PVZUseCase interface {
		CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error)
		GetPVZWithReceptions(ctx context.Context, filter dto.ReceptionFilter) (*[]dto.PVZInfo, error)
	}
	ReceptionUseCase interface {
		CreateReception(ctx context.Context, request dto.ReceptionsRequest) (*entity.Reception, error)
		CloseReception(ctx context.Context, id uuid.UUID) (*entity.Reception, error)
	}
	ProductUseCase interface {
		AddProduct(ctx context.Context, product *dto.PostAddProductRequest) (*entity.Product, error)
		DeleteProductLIFO(ctx context.Context, pvzID uuid.UUID) error
	}
)
