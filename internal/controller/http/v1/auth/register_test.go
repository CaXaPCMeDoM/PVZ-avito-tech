package auth_test

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/controller/http/v1/auth"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/logger"
	authUC "PVZ-avito-tech/internal/usecase/auth"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	loggerMock := logger.NewMock()

	validID := uuid.New()

	tests := []struct {
		name           string
		request        interface{}
		mockAuthSetup  func(*MockAuthUC)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid registration",
			request: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     entity.UserRoleEmployee,
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Register", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Email == "test@example.com" && u.Password == "password123" && u.Role == entity.UserRoleEmployee
				})).Return(authUC.RegisterResponse{
					Id:    validID,
					Email: "test@example.com",
					Role:  entity.UserRoleEmployee,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   validID.String(),
		},
		{
			name:    "invalid request body",
			request: "invalid",
			mockAuthSetup: func(mockAuth *MockAuthUC) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
		{
			name: "invalid role",
			request: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "invalid_role",
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid role",
		},
		{
			name: "user already exists",
			request: dto.RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
				Role:     entity.UserRoleEmployee,
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Register", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Email == "existing@example.com"
				})).Return(authUC.RegisterResponse{}, entity.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "invalid password",
			request: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "weak",
				Role:     entity.UserRoleEmployee,
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Register", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Email == "test@example.com" && u.Password == "weak"
				})).Return(authUC.RegisterResponse{}, entity.ErrInvalidPassword)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   entity.ErrInvalidPassword.Error(),
		},
		{
			name: "password too long",
			request: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "verylongpassword",
				Role:     entity.UserRoleEmployee,
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Register", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Email == "test@example.com" && u.Password == "verylongpassword"
				})).Return(authUC.RegisterResponse{}, entity.ErrPasswordTooLong)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   entity.ErrPasswordTooLong.Error(),
		},
		{
			name: "password hashing error",
			request: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     entity.UserRoleEmployee,
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Register", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Email == "test@example.com" && u.Password == "password123"
				})).Return(authUC.RegisterResponse{}, entity.ErrPasswordHashing)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   entity.ErrPasswordHashing.Error(),
		},
		{
			name: "internal error",
			request: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     entity.UserRoleEmployee,
			},
			mockAuthSetup: func(mockAuth *MockAuthUC) {
				mockAuth.On("Register", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Email == "test@example.com" && u.Password == "password123"
				})).Return(authUC.RegisterResponse{}, entity.ErrInternal)
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
			c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Register(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}

			mockAuth.AssertExpectations(t)
		})
	}
}
