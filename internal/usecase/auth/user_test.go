package auth_test

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/usecase/auth"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, u *entity.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Verify(hashedPassword, inputPassword string) error {
	args := m.Called(hashedPassword, inputPassword)
	return args.Error(0)
}

func TestUserUsecase_Register(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		user          *entity.User
		mockSetup     func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher)
		expectedResp  auth.RegisterResponse
		expectedError error
	}{
		{
			name: "successful registration",
			user: &entity.User{
				Email:    "test@example.com",
				Password: "password123",
				Role:     entity.UserRoleEmployee,
			},
			mockSetup: func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher) {
				mockHasher.On("Hash", "password123").Return("hashed_password", nil)
				mockRepo.On("Create", ctx, mock.MatchedBy(func(u *entity.User) bool {
					return u.Email == "test@example.com" && u.Password == "hashed_password" && u.Role == entity.UserRoleEmployee
				})).Return(nil)
			},
			expectedResp: auth.RegisterResponse{
				Email: "test@example.com",
				Role:  entity.UserRoleEmployee,
			},
			expectedError: nil,
		},
		{
			name: "hashing error",
			user: &entity.User{
				Email:    "test@example.com",
				Password: "password123",
				Role:     entity.UserRoleEmployee,
			},
			mockSetup: func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher) {
				mockHasher.On("Hash", "password123").Return("", entity.ErrPasswordHashing)
			},
			expectedResp:  auth.RegisterResponse{},
			expectedError: entity.ErrPasswordHashing,
		},
		{
			name: "user already exists",
			user: &entity.User{
				Email:    "existing@example.com",
				Password: "password123",
				Role:     entity.UserRoleEmployee,
			},
			mockSetup: func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher) {
				mockHasher.On("Hash", "password123").Return("hashed_password", nil)
				mockRepo.On("Create", ctx, mock.MatchedBy(func(u *entity.User) bool {
					return u.Email == "existing@example.com" && u.Password == "hashed_password"
				})).Return(entity.ErrUserAlreadyExists)
			},
			expectedResp:  auth.RegisterResponse{},
			expectedError: entity.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockHasher := new(MockPasswordHasher)
			uc := auth.NewUserUsecase(mockRepo, mockHasher)

			tt.mockSetup(mockRepo, mockHasher)

			resp, err := uc.Register(ctx, tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResp.Email, resp.Email)
			assert.Equal(t, tt.expectedResp.Role, resp.Role)

			mockRepo.AssertExpectations(t)
			mockHasher.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_Login(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher)
		expectedResp  auth.LoginResponse
		expectedError error
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher) {
				mockRepo.On("GetByEmail", ctx, "test@example.com").Return(&entity.User{
					Email:    "test@example.com",
					Password: "hashed_password",
					Role:     entity.UserRoleModerator,
				}, nil)
				mockHasher.On("Verify", "hashed_password", "password123").Return(nil)
			},
			expectedResp: auth.LoginResponse{
				Role: entity.UserRoleModerator,
			},
			expectedError: nil,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			mockSetup: func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher) {
				mockRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(nil, entity.ErrUserNotFound)
			},
			expectedResp:  auth.LoginResponse{},
			expectedError: entity.ErrUserNotFound,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockSetup: func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher) {
				mockRepo.On("GetByEmail", ctx, "test@example.com").Return(&entity.User{
					Email:    "test@example.com",
					Password: "hashed_password",
					Role:     entity.UserRoleModerator,
				}, nil)
				mockHasher.On("Verify", "hashed_password", "wrongpassword").Return(entity.ErrInvalidPassword)
			},
			expectedResp:  auth.LoginResponse{},
			expectedError: entity.ErrInvalidPassword,
		},
		{
			name:     "internal error",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(mockRepo *MockUserRepo, mockHasher *MockPasswordHasher) {
				mockRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, entity.ErrInternal)
			},
			expectedResp:  auth.LoginResponse{},
			expectedError: entity.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockHasher := new(MockPasswordHasher)
			uc := auth.NewUserUsecase(mockRepo, mockHasher)

			tt.mockSetup(mockRepo, mockHasher)

			resp, err := uc.Login(ctx, tt.email, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResp, resp)

			mockRepo.AssertExpectations(t)
			mockHasher.AssertExpectations(t)
		})
	}
}
