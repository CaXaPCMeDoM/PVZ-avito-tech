package auth_test

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/controller/http/v1/auth"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/logger"
	authUC "PVZ-avito-tech/internal/usecase/auth"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUC struct {
	mock.Mock
}

func (m *MockAuthUC) Register(ctx context.Context, u *entity.User) (authUC.RegisterResponse, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(authUC.RegisterResponse), args.Error(1)
}

func (m *MockAuthUC) Login(ctx context.Context, email string, rawPassword string) (authUC.LoginResponse, error) {
	args := m.Called(ctx, email, rawPassword)
	return args.Get(0).(authUC.LoginResponse), args.Error(1)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	loggerMock := logger.NewMock()

	tests := []struct {
		name           string
		request        interface{}
		mockAuthSetup  func(*MockAuthUC)
		mockDummySetup func(*MockDummyUC)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid login",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Login", mock.Anything, "test@example.com", "password123").
					Return(authUC.LoginResponse{Role: entity.UserRoleModerator}, nil)
			},
			mockDummySetup: func(mockDummy *MockDummyUC) {
				mockDummy.On("GenerateDummyToken", entity.UserRoleModerator).
					Return("test-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"test-token"}`,
		},
		{
			name:    "invalid request body",
			request: "invalid",
			mockAuthSetup: func(mockAuth *MockAuthUC) {
			},
			mockDummySetup: func(mockDummy *MockDummyUC) {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid request body",
		},
		{
			name: "user not found",
			request: dto.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Login", mock.Anything, "nonexistent@example.com", "password123").
					Return(authUC.LoginResponse{}, entity.ErrUserNotFound)
			},
			mockDummySetup: func(mockDummy *MockDummyUC) {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   entity.ErrUserNotFound.Error(),
		},
		{
			name: "invalid password",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Login", mock.Anything, "test@example.com", "wrongpassword").
					Return(authUC.LoginResponse{}, entity.ErrInvalidPassword)
			},
			mockDummySetup: func(mockDummy *MockDummyUC) {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   entity.ErrInvalidPassword.Error(),
		},
		{
			name: "internal error",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Login", mock.Anything, "test@example.com", "password123").
					Return(authUC.LoginResponse{}, entity.ErrInternal)
			},
			mockDummySetup: func(mockDummy *MockDummyUC) {
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   entity.ErrInternal.Error(),
		},
		{
			name: "token generation failure",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Login", mock.Anything, "test@example.com", "password123").
					Return(authUC.LoginResponse{Role: entity.UserRoleModerator}, nil)
			},
			mockDummySetup: func(mockDummy *MockDummyUC) {
				mockDummy.On("GenerateDummyToken", entity.UserRoleModerator).
					Return("", errors.New("token generation failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   entity.ErrInternal.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(MockAuthUC)
			mockDummy := new(MockDummyUC)

			if tt.mockAuthSetup != nil {
				tt.mockAuthSetup(mockAuth)
			}
			if tt.mockDummySetup != nil {
				tt.mockDummySetup(mockDummy)
			}

			router := gin.New()
			handler := auth.NewAuthRoutes(
				router.Group("/"),
				mockDummy,
				mockAuth,
				loggerMock,
			)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.request)
			c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockAuth.AssertExpectations(t)
			mockDummy.AssertExpectations(t)
		})
	}
}
