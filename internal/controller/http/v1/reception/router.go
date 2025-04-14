package reception

import (
	"PVZ-avito-tech/internal/controller/http/middleware"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/auth"
	"PVZ-avito-tech/internal/pkg/logger"
	"PVZ-avito-tech/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	logger      logger.Interface
	receptionUC usecase.ReceptionUseCase
}

func NewAuthRoutes(
	apiV1Group *gin.RouterGroup,
	logger logger.Interface,
	reception usecase.ReceptionUseCase,
	jwtService auth.TokenService,
) *Routes {
	au := &Routes{
		logger:      logger,
		receptionUC: reception,
	}

	authGroup := apiV1Group.Group("/receptions").
		Use(middleware.AuthMiddleware(jwtService, logger))
	{
		authGroup.POST("", middleware.RequireRole(entity.UserRoleEmployee), au.CreateReception)
	}

	return au
}
