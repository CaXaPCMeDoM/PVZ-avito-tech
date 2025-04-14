package dto

import (
	"PVZ-avito-tech/internal/entity"
	_ "github.com/creasty/defaults"
	"github.com/google/uuid"
	"time"
)

type PVZInfo struct {
	PVZ        PVZWithReceptions `json:"pvz"`
	Receptions []ReceptionGroup  `json:"receptions"`
}

type PVZWithReceptions struct {
	ID               uuid.UUID   `json:"id"`
	RegistrationDate time.Time   `json:"registrationDate"`
	City             entity.City `json:"city"`
}

type ReceptionGroup struct {
	Reception ReceptionWithProducts `json:"reception"`
	Products  []ProductDTO          `json:"products"`
}

type ReceptionWithProducts struct {
	ID       uuid.UUID               `json:"id"`
	DateTime time.Time               `json:"dateTime"`
	PVZID    uuid.UUID               `json:"pvzId"`
	Status   entity.ReceptionsStatus `json:"status"`
}

type ProductDTO struct {
	ID          uuid.UUID   `json:"id"`
	DateTime    time.Time   `json:"dateTime"`
	Type        entity.City `json:"type"`
	ReceptionID uuid.UUID   `json:"receptionId"`
}

type CreatePVZRequest struct {
	City             entity.City `json:"city"`
	Id               *uuid.UUID  `json:"id,omitempty"`
	RegistrationDate *time.Time  `json:"registrationDate,omitempty"`
}

type ReceptionFilter struct {
	Page      int       `form:"page" json:"page" binding:"omitempty,min=1" default:"1"`
	Limit     int       `form:"limit" json:"limit" binding:"omitempty,min=1,max=30" default:"10"`
	StartDate time.Time `form:"startDate" json:"startDate" binding:"omitempty,datetime"`
	EndDate   time.Time `form:"endDate" json:"endDate" binding:"omitempty,datetime"`
}

type Option func(*ReceptionFilter)

func WithPaginationDefaults() Option {
	return func(f *ReceptionFilter) {
		if f.Page < 1 {
			f.Page = 1
		}
		if f.Limit < 1 || f.Limit > 30 {
			f.Limit = 10
		}
	}
}

func (f *ReceptionFilter) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(f)
	}
}
