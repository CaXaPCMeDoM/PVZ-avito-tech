package product_test

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/usecase/product"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) AddProduct(ctx context.Context, pvzID uuid.UUID, productType entity.ProductType) (*entity.Product, error) {
	args := m.Called(ctx, pvzID, productType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepo) DeleteProductLIFO(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

func TestUsecase_AddProduct(t *testing.T) {
	ctx := context.Background()
	pvzID := uuid.New()
	productType := entity.ElectronicsProductType
	now := time.Now()
	receptionID := uuid.New()

	tests := []struct {
		name          string
		request       *dto.PostAddProductRequest
		mockSetup     func(*MockProductRepo)
		expectedResp  *entity.Product
		expectedError error
	}{
		{
			name: "successful product addition",
			request: &dto.PostAddProductRequest{
				PvzID:       pvzID,
				ProductType: productType,
			},
			mockSetup: func(mockRepo *MockProductRepo) {
				mockRepo.On("AddProduct", ctx, pvzID, productType).Return(&entity.Product{
					ID:          uuid.New(),
					DateTime:    now,
					Type:        productType,
					ReceptionID: receptionID,
				}, nil)
			},
			expectedResp: &entity.Product{
				Type:        productType,
				ReceptionID: receptionID,
			},
			expectedError: nil,
		},
		{
			name: "error during product addition",
			request: &dto.PostAddProductRequest{
				PvzID:       pvzID,
				ProductType: productType,
			},
			mockSetup: func(mockRepo *MockProductRepo) {
				mockRepo.On("AddProduct", ctx, pvzID, productType).Return(nil, errors.New("failed to add product"))
			},
			expectedResp:  nil,
			expectedError: errors.New("failed to add product"),
		},
		{
			name: "invalid product type",
			request: &dto.PostAddProductRequest{
				PvzID:       pvzID,
				ProductType: "invalid-type",
			},
			mockSetup: func(mockRepo *MockProductRepo) {
				mockRepo.On("AddProduct", ctx, pvzID, entity.ProductType("invalid-type")).Return(nil, errors.New("invalid product type"))
			},
			expectedResp:  nil,
			expectedError: errors.New("invalid product type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepo)
			usecase := product.NewProductUsecase(mockRepo)

			tt.mockSetup(mockRepo)

			resp, err := usecase.AddProduct(ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResp.Type, resp.Type)
				assert.Equal(t, tt.expectedResp.ReceptionID, resp.ReceptionID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUsecase_DeleteProductLIFO(t *testing.T) {
	ctx := context.Background()
	pvzID := uuid.New()

	tests := []struct {
		name          string
		pvzID         uuid.UUID
		mockSetup     func(*MockProductRepo)
		expectedError error
	}{
		{
			name:  "successful product deletion",
			pvzID: pvzID,
			mockSetup: func(mockRepo *MockProductRepo) {
				mockRepo.On("DeleteProductLIFO", ctx, pvzID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "error during product deletion",
			pvzID: pvzID,
			mockSetup: func(mockRepo *MockProductRepo) {
				mockRepo.On("DeleteProductLIFO", ctx, pvzID).Return(errors.New("failed to delete product"))
			},
			expectedError: errors.New("failed to delete product"),
		},
		{
			name:  "pvz not found",
			pvzID: pvzID,
			mockSetup: func(mockRepo *MockProductRepo) {
				mockRepo.On("DeleteProductLIFO", ctx, pvzID).Return(errors.New("pvz not found"))
			},
			expectedError: errors.New("pvz not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepo)
			usecase := product.NewProductUsecase(mockRepo)

			tt.mockSetup(mockRepo)

			err := usecase.DeleteProductLIFO(ctx, tt.pvzID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
