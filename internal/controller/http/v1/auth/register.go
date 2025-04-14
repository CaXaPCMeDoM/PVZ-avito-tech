package auth

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	er "PVZ-avito-tech/internal/controller/http/errors"
	"PVZ-avito-tech/internal/controller/http/mapper"
	"PVZ-avito-tech/internal/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Routes) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidRequestBody)
		return
	}

	if !req.Role.IsValidRole() {
		h.logger.Warn(er.ErrInvalidRole, "role: ", req.Role)
		dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidRole)
		return
	}

	registerResp, err := h.userUC.Register(c.Request.Context(), mapper.RegisterRequestToEntityUser(req))

	if err != nil {
		logFields := map[string]interface{}{
			"email":  req.Email,
			"method": "Register",
		}
		switch {
		case errors.Is(err, entity.ErrUserAlreadyExists):
			h.logger.Warn(err.Error(), logFields)
		case errors.Is(err, entity.ErrInvalidPassword):
			h.logger.Warn(err.Error(), logFields)
			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, entity.ErrPasswordTooLong):
			h.logger.Warn(err.Error(), logFields)
			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, entity.ErrPasswordHashing):
			h.logger.Warn(err.Error(), logFields)
			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			h.logger.Error("unexpected error", logFields)
			dto.ErrorResponse(c, http.StatusInternalServerError, entity.ErrInternal.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, registerResp)
}
