package password_test

import (
	"PVZ-avito-tech/config"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/infrastructure/security/password"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptHasher_Hash(t *testing.T) {

	cfg := &config.Config{
		Security: config.Security{
			PasswordCost: 4,
		},
	}
	hasher := password.NewBcryptHasher(cfg)

	tests := []struct {
		name          string
		password      string
		expectError   bool
		expectedError error
	}{
		{
			name:        "valid password",
			password:    "password123",
			expectError: false,
		},
		{
			name:          "empty password",
			password:      "",
			expectError:   true,
			expectedError: entity.ErrInvalidPassword,
		},
		{
			name:          "password too long",
			password:      strings.Repeat("a", 73),
			expectError:   true,
			expectedError: entity.ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := hasher.Hash(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				assert.True(t, strings.HasPrefix(hashedPassword, "$2a$"))
			}
		})
	}
}

func TestBcryptHasher_Verify(t *testing.T) {
	cfg := &config.Config{
		Security: config.Security{
			PasswordCost: 4,
		},
	}
	hasher := password.NewBcryptHasher(cfg)

	validPassword := "password123"
	hashedPassword, err := hasher.Hash(validPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	tests := []struct {
		name           string
		hashedPassword string
		inputPassword  string
		expectError    bool
		expectedError  error
	}{
		{
			name:           "valid password",
			hashedPassword: hashedPassword,
			inputPassword:  validPassword,
			expectError:    false,
		},
		{
			name:           "invalid password",
			hashedPassword: hashedPassword,
			inputPassword:  "wrongpassword",
			expectError:    true,
			expectedError:  entity.ErrPasswordVerify,
		},
		{
			name:           "invalid hash",
			hashedPassword: "invalid_hash",
			inputPassword:  validPassword,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.Verify(tt.hashedPassword, tt.inputPassword)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
