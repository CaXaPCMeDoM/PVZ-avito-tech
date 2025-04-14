package dto

import (
	"github.com/google/uuid"
)

type ReceptionsRequest struct {
	PvzId uuid.UUID `json:"pvzId"`
}
