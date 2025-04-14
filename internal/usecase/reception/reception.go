package reception

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/infrastructure/repo"
	"context"
	"github.com/google/uuid"
)

type UseCase struct {
	receptionRepo repo.ReceptionRepo
}

func NewUseCase(
	receptionRepo repo.ReceptionRepo,
) *UseCase {
	return &UseCase{
		receptionRepo: receptionRepo,
	}
}

func (uc *UseCase) CreateReception(ctx context.Context, request dto.ReceptionsRequest) (*entity.Reception, error) {
	return uc.receptionRepo.CreateReception(ctx, request.PvzId)
}

func (uc *UseCase) CloseReception(ctx context.Context, id uuid.UUID) (*entity.Reception, error) {
	return uc.receptionRepo.CloseActiveReception(ctx, id)
}
