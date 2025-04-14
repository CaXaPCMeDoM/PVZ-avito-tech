package products

import (
	"PVZ-avito-tech/internal/controller/http/middleware"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/auth"
	"PVZ-avito-tech/internal/pkg/logger"
	"PVZ-avito-tech/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	logger    logger.Interface
	productUC usecase.ProductUseCase
}

func NewAuthRoutes(
	apiV1Group *gin.RouterGroup,
	productUC usecase.ProductUseCase,
	logger logger.Interface,
	jwtService auth.TokenService,
) *Routes {
	au := &Routes{
		logger:    logger,
		productUC: productUC,
	}

	authGroup := apiV1Group.Group("/products").
		Use(middleware.AuthMiddleware(jwtService, logger))
	{
		authGroup.POST("", middleware.RequireRole(entity.UserRoleEmployee), au.AddProduct)
	}

	return au
}
