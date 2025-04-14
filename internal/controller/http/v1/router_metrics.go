package v1

import (
	"PVZ-avito-tech/internal/controller/http/middleware"
	"PVZ-avito-tech/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouterMetrics(
	l logger.Interface,
) *gin.Engine {
	router := gin.New()

	router.Use(
		middleware.Logger(l),
		middleware.PrometheusMiddleware(),
	)

	apiV1 := router.Group("/metrics")
	{
		apiV1.GET("", gin.WrapH(promhttp.Handler()))
	}

	return router
}
