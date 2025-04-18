package dummy_test

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/auth"
	"PVZ-avito-tech/internal/usecase/dummy"
	"errors"
	"testing"

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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Claims), args.Error(1)
}

func TestAuthUseCase_GenerateDummyToken(t *testing.T) {
	tests := []struct {
		name          string
		role          entity.UserRole
		mockSetup     func(*MockTokenService)
		expectedToken string
		expectedError error
	}{
		{
			name: "successful token generation for employee",
			role: entity.UserRoleEmployee,
			mockSetup: func(mockJWT *MockTokenService) {
				mockJWT.On("Generate", entity.UserRoleEmployee).Return("employee-token", nil)
			},
			expectedToken: "employee-token",
			expectedError: nil,
		},
		{
			name: "successful token generation for moderator",
			role: entity.UserRoleModerator,
			mockSetup: func(mockJWT *MockTokenService) {
				mockJWT.On("Generate", entity.UserRoleModerator).Return("moderator-token", nil)
			},
			expectedToken: "moderator-token",
			expectedError: nil,
		},
		{
			name: "error during token generation",
			role: entity.UserRoleEmployee,
			mockSetup: func(mockJWT *MockTokenService) {
				mockJWT.On("Generate", entity.UserRoleEmployee).Return("", errors.New("token generation failed"))
			},
			expectedToken: "",
			expectedError: errors.New("token generation failed"),
		},
		{
			name: "invalid role",
			role: "invalid-role",
			mockSetup: func(mockJWT *MockTokenService) {
				mockJWT.On("Generate", entity.UserRole("invalid-role")).Return("", errors.New("invalid role"))
			},
			expectedToken: "",
			expectedError: errors.New("invalid role"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWT := new(MockTokenService)
			usecase := dummy.NewDummyAuthUseCase(mockJWT)

			tt.mockSetup(mockJWT)

			token, err := usecase.GenerateDummyToken(tt.role)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedToken, token)

			mockJWT.AssertExpectations(t)
		})
	}
}
