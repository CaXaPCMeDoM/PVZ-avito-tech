package jwt_test

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/auth"
	jwtpkg "PVZ-avito-tech/internal/pkg/auth/jwt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name        string
		secretKey   []byte
		expectError bool
	}{
		{
			name:        "valid secret",
			secretKey:   []byte("valid-secret-key"),
			expectError: false,
		},
		{
			name:        "empty secret",
			secretKey:   []byte{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := jwtpkg.NewService(tt.secretKey)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, service)
				assert.Equal(t, jwtpkg.ErrEmptySecret, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	secretKey := []byte("test-secret-key")
	service, err := jwtpkg.NewService(secretKey)
	require.NoError(t, err)
	require.NotNil(t, service)

	tests := []struct {
		name      string
		role      entity.UserRole
		expectErr bool
	}{
		{
			name:      "employee role",
			role:      entity.UserRoleEmployee,
			expectErr: false,
		},
		{
			name:      "moderator role",
			role:      entity.UserRoleModerator,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Generate(tt.role)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				claims := &auth.Claims{}
				parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
					return secretKey, nil
				})

				assert.NoError(t, err)
				assert.True(t, parsedToken.Valid)
				assert.Equal(t, tt.role, claims.Role)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	secretKey := []byte("test-secret-key")
	service, err := jwtpkg.NewService(secretKey)
	require.NoError(t, err)
	require.NotNil(t, service)

	validRole := entity.UserRoleModerator
	validToken, err := service.Generate(validRole)
	require.NoError(t, err)
	require.NotEmpty(t, validToken)

	invalidSecretService, err := jwtpkg.NewService([]byte("different-secret"))
	require.NoError(t, err)
	invalidToken, err := invalidSecretService.Generate(validRole)
	require.NoError(t, err)

	expiredClaims := auth.Claims{
		Role: validRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, err := expiredToken.SignedString(secretKey)
	require.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		expectErr bool
		role      entity.UserRole
	}{
		{
			name:      "valid token",
			token:     validToken,
			expectErr: false,
			role:      validRole,
		},
		{
			name:      "empty token",
			token:     "",
			expectErr: true,
		},
		{
			name:      "invalid token format",
			token:     "invalid.token.format",
			expectErr: true,
		},
		{
			name:      "token with invalid signature",
			token:     invalidToken,
			expectErr: true,
		},
		{
			name:      "expired token",
			token:     expiredTokenString,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.Validate(tt.token)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
				assert.Equal(t, auth.ErrInvalidToken, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.role, claims.Role)
			}
		})
	}
}
