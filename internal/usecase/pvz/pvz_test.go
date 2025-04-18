package pvz_test

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/logger"
	"PVZ-avito-tech/internal/usecase/pvz"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockPVZRepo struct {
	mock.Mock
}

func (m *MockPVZRepo) Create(ctx context.Context, pvz *entity.PVZ) error {
	args := m.Called(ctx, pvz)
	return args.Error(0)
}

func (m *MockPVZRepo) GetPVZWithReceptions(ctx context.Context, filter dto.ReceptionFilter) (*[]dto.PVZInfo, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.PVZInfo), args.Error(1)
}

type MockReceptionRepo struct {
	mock.Mock
}

func (m *MockReceptionRepo) CreateReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error) {
	args := m.Called(ctx, pvzID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Reception), args.Error(1)
}

func (m *MockReceptionRepo) CloseActiveReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error) {
	args := m.Called(ctx, pvzID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Reception), args.Error(1)
}

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

func TestUseCase_CreatePVZ(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	now := time.Now()

	tests := []struct {
		name          string
		pvz           *entity.PVZ
		mockSetup     func(*MockPVZRepo, *MockReceptionRepo, *MockProductRepo)
		expectedResp  *entity.PVZ
		expectedError error
	}{
		{
			name: "successful pvz creation",
			pvz: &entity.PVZ{
				ID:               &id,
				City:             entity.CityMoscow,
				RegistrationDate: &now,
			},
			mockSetup: func(mockPVZRepo *MockPVZRepo, mockReceptionRepo *MockReceptionRepo, mockProductRepo *MockProductRepo) {
				mockPVZRepo.On("Create", ctx, mock.MatchedBy(func(p *entity.PVZ) bool {
					return p.ID == &id && p.City == entity.CityMoscow && p.RegistrationDate == &now
				})).Return(nil)
			},
			expectedResp: &entity.PVZ{
				ID:               &id,
				City:             entity.CityMoscow,
				RegistrationDate: &now,
			},
			expectedError: nil,
		},
		{
			name: "error during pvz creation",
			pvz: &entity.PVZ{
				ID:               &id,
				City:             entity.CityMoscow,
				RegistrationDate: &now,
			},
			mockSetup: func(mockPVZRepo *MockPVZRepo, mockReceptionRepo *MockReceptionRepo, mockProductRepo *MockProductRepo) {
				mockPVZRepo.On("Create", ctx, mock.MatchedBy(func(p *entity.PVZ) bool {
					return p.ID == &id && p.City == entity.CityMoscow && p.RegistrationDate == &now
				})).Return(errors.New("failed to create pvz"))
			},
			expectedResp:  nil,
			expectedError: errors.New("failed to create pvz"),
		},
		{
			name: "invalid city",
			pvz: &entity.PVZ{
				ID:               &id,
				City:             "Invalid City",
				RegistrationDate: &now,
			},
			mockSetup: func(mockPVZRepo *MockPVZRepo, mockReceptionRepo *MockReceptionRepo, mockProductRepo *MockProductRepo) {
				mockPVZRepo.On("Create", ctx, mock.MatchedBy(func(p *entity.PVZ) bool {
					return p.ID == &id && p.City == "Invalid City" && p.RegistrationDate == &now
				})).Return(errors.New("invalid city"))
			},
			expectedResp:  nil,
			expectedError: errors.New("invalid city"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZRepo := new(MockPVZRepo)
			mockReceptionRepo := new(MockReceptionRepo)
			mockProductRepo := new(MockProductRepo)
			loggerMock := logger.NewMock()

			usecase := pvz.NewPVZUseCase(mockPVZRepo, mockReceptionRepo, mockProductRepo, loggerMock)

			tt.mockSetup(mockPVZRepo, mockReceptionRepo, mockProductRepo)

			resp, err := usecase.CreatePVZ(ctx, tt.pvz)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResp.ID, resp.ID)
				assert.Equal(t, tt.expectedResp.City, resp.City)
				assert.Equal(t, tt.expectedResp.RegistrationDate, resp.RegistrationDate)
			}

			mockPVZRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
			mockProductRepo.AssertExpectations(t)
		})
	}
}

func TestUseCase_GetPVZWithReceptions(t *testing.T) {
	ctx := context.Background()
	filter := dto.ReceptionFilter{
		Page:  1,
		Limit: 10,
	}

	pvzID := uuid.New()
	receptionID := uuid.New()
	productID := uuid.New()
	now := time.Now()

	tests := []struct {
		name          string
		filter        dto.ReceptionFilter
		mockSetup     func(*MockPVZRepo, *MockReceptionRepo, *MockProductRepo)
		expectedResp  *[]dto.PVZInfo
		expectedError error
	}{
		{
			name:   "successful get pvz with receptions",
			filter: filter,
			mockSetup: func(mockPVZRepo *MockPVZRepo, mockReceptionRepo *MockReceptionRepo, mockProductRepo *MockProductRepo) {
				mockPVZRepo.On("GetPVZWithReceptions", ctx, filter).Return(&[]dto.PVZInfo{
					{
						PVZ: dto.PVZWithReceptions{
							ID:               pvzID,
							RegistrationDate: now,
							City:             entity.CityMoscow,
						},
						Receptions: []*dto.ReceptionGroup{
							{
								Reception: dto.ReceptionWithProducts{
									ID:       receptionID,
									DateTime: now,
									PVZID:    pvzID,
									Status:   entity.InProgressStatus,
								},
								Products: []dto.ProductDTO{
									{
										ID:          productID,
										DateTime:    now,
										Type:        entity.ElectronicsProductType,
										ReceptionID: receptionID,
									},
								},
							},
						},
					},
				}, nil)
			},
			expectedResp: &[]dto.PVZInfo{
				{
					PVZ: dto.PVZWithReceptions{
						ID:               pvzID,
						RegistrationDate: now,
						City:             entity.CityMoscow,
					},
					Receptions: []*dto.ReceptionGroup{
						{
							Reception: dto.ReceptionWithProducts{
								ID:       receptionID,
								DateTime: now,
								PVZID:    pvzID,
								Status:   entity.InProgressStatus,
							},
							Products: []dto.ProductDTO{
								{
									ID:          productID,
									DateTime:    now,
									Type:        entity.ElectronicsProductType,
									ReceptionID: receptionID,
								},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:   "error during get pvz with receptions",
			filter: filter,
			mockSetup: func(mockPVZRepo *MockPVZRepo, mockReceptionRepo *MockReceptionRepo, mockProductRepo *MockProductRepo) {
				mockPVZRepo.On("GetPVZWithReceptions", ctx, filter).Return(nil, errors.New("failed to get pvz with receptions"))
			},
			expectedResp:  nil,
			expectedError: errors.New("failed to get pvz with receptions"),
		},
		{
			name:   "no pvz found",
			filter: filter,
			mockSetup: func(mockPVZRepo *MockPVZRepo, mockReceptionRepo *MockReceptionRepo, mockProductRepo *MockProductRepo) {
				mockPVZRepo.On("GetPVZWithReceptions", ctx, filter).Return(&[]dto.PVZInfo{}, nil)
			},
			expectedResp:  &[]dto.PVZInfo{},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZRepo := new(MockPVZRepo)
			mockReceptionRepo := new(MockReceptionRepo)
			mockProductRepo := new(MockProductRepo)
			loggerMock := logger.NewMock()

			usecase := pvz.NewPVZUseCase(mockPVZRepo, mockReceptionRepo, mockProductRepo, loggerMock)

			tt.mockSetup(mockPVZRepo, mockReceptionRepo, mockProductRepo)

			resp, err := usecase.GetPVZWithReceptions(ctx, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}

			mockPVZRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
			mockProductRepo.AssertExpectations(t)
		})
	}
}
