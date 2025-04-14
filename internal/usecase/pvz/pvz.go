package pvz

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/infrastructure/repo"
	"PVZ-avito-tech/internal/pkg/logger"
	"context"
)

type UseCase struct {
	pvzRepo       repo.PVZRepo
	receptionRepo repo.ReceptionRepo
	productRepo   repo.ProductRepo
	log           logger.Interface
}

func NewPVZUseCase(
	pvzRepo repo.PVZRepo,
	receptionRepo repo.ReceptionRepo,
	productRepo repo.ProductRepo,
	log logger.Interface,
) *UseCase {
	return &UseCase{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
		log:           log,
	}
}

func (uc *UseCase) CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error) {
	err := uc.pvzRepo.Create(ctx, pvz)
	if err != nil {
		return nil, err
	}

	return pvz, err
}

func (uc *UseCase) GetPVZWithReceptions(ctx context.Context, filter dto.ReceptionFilter) (*[]dto.PVZInfo, error) {
	return uc.pvzRepo.GetPVZWithReceptions(ctx, filter)
}
