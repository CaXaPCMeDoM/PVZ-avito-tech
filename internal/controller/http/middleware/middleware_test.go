package middleware_test

import (
	"PVZ-avito-tech/internal/controller/http/middleware"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/auth"
	"PVZ-avito-tech/internal/pkg/logger"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) Generate(role entity.UserRole) (string, error) {
	args := m.Called(role)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) Validate(token string) (*auth.Claims, error) {
	args := m.Called(token)
	return args.Get(0).(*auth.Claims), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	loggerMock := logger.NewMock()
	tokenService := new(MockTokenService)

	tests := []struct {
		name           string
		authHeader     string
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "no authorization header",
			authHeader:     "",
			expectedStatus: http.StatusForbidden,
			expectedError:  "authorization header is required",
		},
		{
			name:           "invalid header format",
			authHeader:     "Basic token",
			expectedStatus: http.StatusForbidden,
			expectedError:  "invalid authorization header format",
		},
		{
			name:       "invalid token",
			authHeader: "Bearer bad_token",
			mockSetup: func() {
				tokenService.On("Validate", "bad_token").Return(
					&auth.Claims{},
					errors.New("invalid token"),
				)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "invalid token",
		},
		{
			name:       "valid token with invalid role",
			authHeader: "Bearer valid_token",
			mockSetup: func() {
				tokenService.On("Validate", "valid_token").Return(
					&auth.Claims{Role: "invalid_role"},
					nil,
				)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "invalid role in token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			r := gin.New()
			r.Use(middleware.AuthMiddleware(tokenService, loggerMock))
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			tokenService.AssertExpectations(t)
		})
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		contextRole    interface{}
		requiredRoles  []entity.UserRole
		expectedStatus int
	}{
		{
			name:           "role not in context",
			contextRole:    nil,
			requiredRoles:  []entity.UserRole{entity.UserRoleModerator},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "invalid role type",
			contextRole:    "invalid_type",
			requiredRoles:  []entity.UserRole{entity.UserRoleModerator},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "access granted",
			contextRole:    entity.UserRoleModerator,
			requiredRoles:  []entity.UserRole{entity.UserRoleModerator},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "access denied",
			contextRole:    entity.UserRoleEmployee,
			requiredRoles:  []entity.UserRole{entity.UserRoleModerator},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(func(c *gin.Context) {
				if tt.contextRole != nil {
					c.Set(middleware.UserRoleContextKey, tt.contextRole)
				}
			})
			r.GET("/test", middleware.RequireRole(tt.requiredRoles...), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
