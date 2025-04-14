package v1

import (
	"PVZ-avito-tech/config"
	"PVZ-avito-tech/internal/controller/http/middleware"
	"PVZ-avito-tech/internal/controller/http/v1/auth"
	"PVZ-avito-tech/internal/controller/http/v1/products"
	"PVZ-avito-tech/internal/controller/http/v1/pvz"
	"PVZ-avito-tech/internal/controller/http/v1/reception"
	authPkg "PVZ-avito-tech/internal/pkg/auth"
	"PVZ-avito-tech/internal/pkg/logger"
	"PVZ-avito-tech/internal/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	l logger.Interface,
	authUC usecase.Auth,
	dummyAuthUC usecase.DummyLogin,
	receptionUC usecase.ReceptionUseCase,
	pvzUC usecase.PVZUseCase,
	productUC usecase.ProductUseCase,
	jwtService authPkg.TokenService,
) *gin.Engine {
	router := gin.New()

	router.Use(
		middleware.Logger(l),
		middleware.PrometheusMiddleware(),
	)

	apiV1 := router.Group("")
	{
		auth.NewAuthRoutes(
			apiV1,
			dummyAuthUC,
			authUC,
			l,
		)

		pvz.NewAuthRoutes(
			apiV1,
			l,
			pvzUC,
			receptionUC,
			productUC,
			jwtService,
		)

		reception.NewAuthRoutes(
			apiV1,
			l,
			receptionUC,
			jwtService,
		)

		products.NewAuthRoutes(
			apiV1,
			productUC,
			l,
			jwtService,
		)
	}

	return router
}
