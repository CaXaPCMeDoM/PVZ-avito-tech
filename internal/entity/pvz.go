package entity

import (
	"github.com/google/uuid"
	"time"
)

type PVZ struct {
	ID               *uuid.UUID `json:"id"`
	City             City       `json:"city"`
	RegistrationDate *time.Time `json:"registrationDate"`
}
