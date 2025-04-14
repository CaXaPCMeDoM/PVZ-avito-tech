package product

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/infrastructure/repo"
	"context"
	"github.com/google/uuid"
)

type Usecase struct {
	repo repo.ProductRepo
}

func NewProductUsecase(repo repo.ProductRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) AddProduct(ctx context.Context, product *dto.PostAddProductRequest) (*entity.Product, error) {
	return uc.repo.AddProduct(ctx, product.PvzID, product.ProductType)
}

func (uc *Usecase) DeleteProductLIFO(ctx context.Context, pvzID uuid.UUID) error {
	return uc.repo.DeleteProductLIFO(ctx, pvzID)
}
