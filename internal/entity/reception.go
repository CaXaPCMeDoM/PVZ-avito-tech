package entity

import (
	"github.com/google/uuid"
	"time"
)

type Reception struct {
	ID       uuid.UUID        `json:"id"`
	DateTime time.Time        `json:"dateTime"`
	PVZID    uuid.UUID        `json:"pvzId"`
	Status   ReceptionsStatus `json:"status"`
}
