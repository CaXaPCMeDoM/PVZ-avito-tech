package auth_test

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/controller/http/errors"
	"PVZ-avito-tech/internal/controller/http/v1/auth"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/logger"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDummyUC struct {
	mock.Mock
}

func (m *MockDummyUC) GenerateDummyToken(role entity.UserRole) (string, error) {
	args := m.Called(role)
	return args.String(0), args.Error(1)
}

func TestDummyLogin(t *testing.T) {
	loggerMock := logger.NewMock()
	mockDummyUC := new(MockDummyUC)
	handler := auth.NewAuthRoutes(
		gin.New().Group("/"),
		mockDummyUC,
		nil,
		loggerMock,
	)

	tests := []struct {
		name           string
		request        interface{}
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid request body",
			request:        "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.ErrInvalidRequestBody,
		},
		{
			name: "invalid role",
			request: dto.DummyLoginRequest{
				Role: "invalid_role",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.ErrInvalidRole,
		},
		{
			name: "success",
			request: dto.DummyLoginRequest{
				Role: "moderator",
			},
			mockSetup: func() {
				mockDummyUC.On("GenerateDummyToken", entity.UserRoleModerator).Return("test-token", nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.request)
			c.Request = httptest.NewRequest("POST", "/dummyLogin", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.DummyLogin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockDummyUC.AssertExpectations(t)
		})
	}
}
