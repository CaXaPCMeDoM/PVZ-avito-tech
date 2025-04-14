package middleware

import (
	"PVZ-avito-tech/internal/pkg/metrics"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		metrics.RequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		metrics.RequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}
