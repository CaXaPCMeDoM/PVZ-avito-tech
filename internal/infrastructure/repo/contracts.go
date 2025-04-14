package repo

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"context"
	"github.com/google/uuid"
)

type (
	UserRepo interface {
		Create(ctx context.Context, u *entity.User) error
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
	}

	PVZRepo interface {
		Create(ctx context.Context, pvz *entity.PVZ) error
		GetPVZWithReceptions(ctx context.Context, filter dto.ReceptionFilter) (*[]dto.PVZInfo, error)
	}

	ReceptionRepo interface {
		CreateReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error)
		CloseActiveReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error)
	}

	ProductRepo interface {
		AddProduct(ctx context.Context, pvzID uuid.UUID, productType entity.ProductType) (*entity.Product, error)
		DeleteProductLIFO(ctx context.Context, pvzID uuid.UUID) error
	}
)
