package mapper

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
)

func EntityProductToProductResponse(ent *entity.Product) *dto.PostAddProductResponse {
	return &dto.PostAddProductResponse{
		ID:          ent.ID,
		DateTime:    ent.DateTime,
		Type:        ent.Type,
		ReceptionID: ent.ReceptionID,
	}
}
