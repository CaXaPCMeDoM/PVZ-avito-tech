package middleware

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/auth"
	"PVZ-avito-tech/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	UserRoleContextKey  = "userRole"
	BearerSchema        = "Bearer "
)

func AuthMiddleware(jwtService auth.TokenService, logger logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "authorization header is required"})
			return
		}

		if !strings.HasPrefix(authHeader, BearerSchema) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid authorization header format"})
			return
		}

		token := strings.TrimPrefix(authHeader, BearerSchema)
		claims, err := jwtService.Validate(token)
		if err != nil {
			logger.Error("JWT validation failed: %v", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}

		if err := claims.Role.ValidateRole(); err != nil {
			logger.Error("Invalid role in token: %s", claims.Role)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid role in token"})
			return
		}

		c.Set(UserRoleContextKey, claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...entity.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(UserRoleContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "role not found in context"})
			return
		}

		userRole, ok := roleValue.(entity.UserRole)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid role type in context"})
			return
		}

		for _, allowedRole := range roles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden,
			gin.H{"error": "required role: " + strings.Join(userRolesToStrings(roles), ", ")})
	}
}

func userRolesToStrings(roles []entity.UserRole) []string {
	res := make([]string, len(roles))
	for i, r := range roles {
		res[i] = string(r)
	}
	return res
}
