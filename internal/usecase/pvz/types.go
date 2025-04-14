package pvz

import (
	"PVZ-avito-tech/internal/entity"
	"github.com/google/uuid"
	"time"
)

type CreateResponse struct {
	Id               uuid.UUID   `json:"id"`
	City             entity.City `json:"city"`
	RegistrationDate time.Time   `json:"registrationDate"`
}

type WithReceptions struct {
	PVZ        *entity.PVZ
	Receptions []*ReceptionWithProducts
}

type ReceptionWithProducts struct {
	Reception *entity.Reception
	Products  []*entity.Product
}
