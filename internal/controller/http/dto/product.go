package dto

import (
	"PVZ-avito-tech/internal/entity"
	"github.com/google/uuid"
	"time"
)

type PostAddProductRequest struct {
	PvzID       uuid.UUID          `json:"pvzId"`
	ProductType entity.ProductType `json:"type"`
}

type PostAddProductResponse struct {
	ID          uuid.UUID          `json:"id"`
	DateTime    time.Time          `json:"dateTime"`
	Type        entity.ProductType `json:"type"`
	ReceptionID uuid.UUID          `json:"receptionId"`
}
