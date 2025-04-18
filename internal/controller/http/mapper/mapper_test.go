package mapper_test

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/controller/http/mapper"
	"PVZ-avito-tech/internal/entity"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterRequestToEntityUser(t *testing.T) {
	// Arrange
	email := "test@example.com"
	password := "password123"
	role := entity.UserRoleEmployee

	request := dto.RegisterRequest{
		Email:    email,
		Password: password,
		Role:     role,
	}

	// Act
	user := mapper.RegisterRequestToEntityUser(request)

	// Assert
	assert.Equal(t, email, user.Email)
	assert.Equal(t, password, user.Password)
	assert.Equal(t, role, user.Role)
}

func TestDtoPVZToEntityPVZ(t *testing.T) {
	// Arrange
	id := uuid.New()
	city := entity.CityMoscow
	now := time.Now()

	request := dto.CreatePVZRequest{
		City:             city,
		Id:               &id,
		RegistrationDate: &now,
	}

	// Act
	pvz := mapper.DtoPVZToEntityPVZ(request)

	// Assert
	assert.Equal(t, &id, pvz.ID)
	assert.Equal(t, city, pvz.City)
	assert.Equal(t, &now, pvz.RegistrationDate)
}

func TestEntityProductToProductResponse(t *testing.T) {
	// Arrange
	id := uuid.New()
	receptionID := uuid.New()
	now := time.Now()
	productType := entity.ElectronicsProductType

	product := &entity.Product{
		ID:          id,
		DateTime:    now,
		Type:        productType,
		ReceptionID: receptionID,
	}

	// Act
	response := mapper.EntityProductToProductResponse(product)

	// Assert
	assert.Equal(t, id, response.ID)
	assert.Equal(t, now, response.DateTime)
	assert.Equal(t, productType, response.Type)
	assert.Equal(t, receptionID, response.ReceptionID)
}
