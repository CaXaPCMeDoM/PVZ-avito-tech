package auth

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/controller/http/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Routes) DummyLogin(c *gin.Context) {
	var req dto.DummyLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(errors.ErrInvalidRequestBody)
		dto.ErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidRequestBody)
		return
	}

	role := req.Role
	if !role.IsValidRole() {
		h.logger.Warn(errors.ErrInvalidRole, "role: ", req.Role)
		dto.ErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidRole)
		return
	}

	token, err := h.dummyUC.GenerateDummyToken(role)
	if err != nil {
		h.logger.Error(err, errors.ErrTokenGeneration, "http - v1 - auth - handler - DummyLogin")
		dto.ErrorResponse(c, http.StatusInternalServerError, errors.ErrTokenGeneration)
		return
	}

	response := tokenResponse{Token: token}

	c.JSON(http.StatusOK, response)
}
