package pvz

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
	pvzUC       usecase.PVZUseCase
	receptionUC usecase.ReceptionUseCase
	productUC   usecase.ProductUseCase
}

func NewAuthRoutes(
	apiV1Group *gin.RouterGroup,
	logger logger.Interface,
	pvzUC usecase.PVZUseCase,
	receptionUC usecase.ReceptionUseCase,
	productUC usecase.ProductUseCase,
	jwtService auth.TokenService,
) *Routes {
	au := &Routes{
		logger:      logger,
		pvzUC:       pvzUC,
		receptionUC: receptionUC,
		productUC:   productUC,
	}

	authGroup := apiV1Group.Group("/pvz").
		Use(middleware.AuthMiddleware(jwtService, logger))
	{
		authGroup.POST("", middleware.RequireRole(entity.UserRoleModerator), au.CreatePVZ)
		authGroup.GET("", middleware.RequireRole(entity.UserRoleModerator, entity.UserRoleEmployee), au.GetPVZList)
		authGroup.POST("/:pvzId/close_last_reception", middleware.RequireRole(entity.UserRoleEmployee), au.CloseReception)
		authGroup.POST("/:pvzId/delete_last_product", middleware.RequireRole(entity.UserRoleEmployee), au.DeleteLastProduct)
	}

	return au
}
