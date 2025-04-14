package dto_test

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestCreatePVZRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   dto.CreatePVZRequest
		wantErr bool
	}{
		{
			name:    "valid city",
			input:   dto.CreatePVZRequest{City: entity.CityMoscow},
			wantErr: false,
		},
		{
			name:    "invalid city",
			input:   dto.CreatePVZRequest{City: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.City.ValidateCity()
			if (err != nil) != tt.wantErr {
				t.Errorf("City validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReceptionFilter_ApplyOptions(t *testing.T) {
	filter := &dto.ReceptionFilter{
		Page:  0,
		Limit: 35,
	}

	filter.Apply(dto.WithPaginationDefaults())

	if filter.Page != 1 {
		t.Errorf("Expected page 1, got %d", filter.Page)
	}
	if filter.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", filter.Limit)
	}
}

func TestPVZWithReceptions_Structure(t *testing.T) {
	pvz := dto.PVZWithReceptions{
		ID:               uuid.New(),
		City:             entity.CityKazan,
		RegistrationDate: time.Now(),
	}

	if pvz.ID == uuid.Nil {
		t.Error("Expected valid UUID")
	}
	if !pvz.City.IsValidCity() {
		t.Errorf("Invalid city: %s", pvz.City)
	}
}
