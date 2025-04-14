package mapper

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
)

func DtoPVZToEntityPVZ(pvz dto.CreatePVZRequest) *entity.PVZ {
	return &entity.PVZ{
		ID:               pvz.Id,
		City:             pvz.City,
		RegistrationDate: pvz.RegistrationDate,
	}
}
