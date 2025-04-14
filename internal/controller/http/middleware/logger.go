package middleware

import (
	"PVZ-avito-tech/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func buildRequestMessage(c *gin.Context, duration time.Duration) string {
	var result strings.Builder

	result.WriteString(c.ClientIP())
	result.WriteString(" - ")
	result.WriteString(c.Request.Method)
	result.WriteString(" ")
	result.WriteString(c.Request.RequestURI)
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(c.Writer.Status()))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(c.Writer.Size()))
	result.WriteString(" - ")
	result.WriteString(duration.String())

	return result.String()
}

func Logger(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		msg := buildRequestMessage(c, time.Since(start))
		l.Info(msg)
	}
}
