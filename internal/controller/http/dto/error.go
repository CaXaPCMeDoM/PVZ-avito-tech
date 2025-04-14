package dto

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, Error{Message: msg})
}

type Error struct {
	Message string `json:"message"`
}
