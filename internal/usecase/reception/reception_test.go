package reception_test

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/usecase/reception"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

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

func TestUseCase_CreateReception(t *testing.T) {
	ctx := context.Background()
	pvzID := uuid.New()
	receptionID := uuid.New()
	now := time.Now()

	tests := []struct {
		name          string
		request       dto.ReceptionsRequest
		mockSetup     func(*MockReceptionRepo)
		expectedResp  *entity.Reception
		expectedError error
	}{
		{
			name: "successful reception creation",
			request: dto.ReceptionsRequest{
				PvzId: pvzID,
			},
			mockSetup: func(mockRepo *MockReceptionRepo) {
				mockRepo.On("CreateReception", ctx, pvzID).Return(&entity.Reception{
					ID:       receptionID,
					DateTime: now,
					PVZID:    pvzID,
					Status:   entity.InProgressStatus,
				}, nil)
			},
			expectedResp: &entity.Reception{
				ID:       receptionID,
				DateTime: now,
				PVZID:    pvzID,
				Status:   entity.InProgressStatus,
			},
			expectedError: nil,
		},
		{
			name: "error during reception creation",
			request: dto.ReceptionsRequest{
				PvzId: pvzID,
			},
			mockSetup: func(mockRepo *MockReceptionRepo) {
				mockRepo.On("CreateReception", ctx, pvzID).Return(nil, errors.New("failed to create reception"))
			},
			expectedResp:  nil,
			expectedError: errors.New("failed to create reception"),
		},
		{
			name: "pvz not found",
			request: dto.ReceptionsRequest{
				PvzId: pvzID,
			},
			mockSetup: func(mockRepo *MockReceptionRepo) {
				mockRepo.On("CreateReception", ctx, pvzID).Return(nil, errors.New("pvz not found"))
			},
			expectedResp:  nil,
			expectedError: errors.New("pvz not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockReceptionRepo)
			usecase := reception.NewUseCase(mockRepo)

			tt.mockSetup(mockRepo)

			resp, err := usecase.CreateReception(ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResp.ID, resp.ID)
				assert.Equal(t, tt.expectedResp.PVZID, resp.PVZID)
				assert.Equal(t, tt.expectedResp.Status, resp.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUseCase_CloseReception(t *testing.T) {
	ctx := context.Background()
	pvzID := uuid.New()
	receptionID := uuid.New()
	now := time.Now()

	tests := []struct {
		name          string
		pvzID         uuid.UUID
		mockSetup     func(*MockReceptionRepo)
		expectedResp  *entity.Reception
		expectedError error
	}{
		{
			name:  "successful reception closing",
			pvzID: pvzID,
			mockSetup: func(mockRepo *MockReceptionRepo) {
				mockRepo.On("CloseActiveReception", ctx, pvzID).Return(&entity.Reception{
					ID:       receptionID,
					DateTime: now,
					PVZID:    pvzID,
					Status:   entity.CloseStatus,
				}, nil)
			},
			expectedResp: &entity.Reception{
				ID:       receptionID,
				DateTime: now,
				PVZID:    pvzID,
				Status:   entity.CloseStatus,
			},
			expectedError: nil,
		},
		{
			name:  "error during reception closing",
			pvzID: pvzID,
			mockSetup: func(mockRepo *MockReceptionRepo) {
				mockRepo.On("CloseActiveReception", ctx, pvzID).Return(nil, errors.New("failed to close reception"))
			},
			expectedResp:  nil,
			expectedError: errors.New("failed to close reception"),
		},
		{
			name:  "pvz not found",
			pvzID: pvzID,
			mockSetup: func(mockRepo *MockReceptionRepo) {
				mockRepo.On("CloseActiveReception", ctx, pvzID).Return(nil, errors.New("pvz not found"))
			},
			expectedResp:  nil,
			expectedError: errors.New("pvz not found"),
		},
		{
			name:  "no active reception",
			pvzID: pvzID,
			mockSetup: func(mockRepo *MockReceptionRepo) {
				mockRepo.On("CloseActiveReception", ctx, pvzID).Return(nil, errors.New("no active reception"))
			},
			expectedResp:  nil,
			expectedError: errors.New("no active reception"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockReceptionRepo)
			usecase := reception.NewUseCase(mockRepo)

			tt.mockSetup(mockRepo)

			resp, err := usecase.CloseReception(ctx, tt.pvzID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResp.ID, resp.ID)
				assert.Equal(t, tt.expectedResp.PVZID, resp.PVZID)
				assert.Equal(t, tt.expectedResp.Status, resp.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
