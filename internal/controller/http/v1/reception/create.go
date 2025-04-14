package reception

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	er "PVZ-avito-tech/internal/controller/http/errors"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/metrics"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Routes) CreateReception(c *gin.Context) {
	var req dto.ReceptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidRequestBody)
	}

	response, err := h.receptionUC.CreateReception(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrReceptionConflict):
			h.logger.Warn(err.Error())
			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, entity.ErrPVZNotFound):
			h.logger.Warn(err.Error())
			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			h.logger.Error(err.Error())
			dto.ErrorResponse(c, http.StatusInternalServerError, entity.ErrInternal.Error())
		}
		return
	}

	metrics.ReceptionsCreated.Inc()
	c.JSON(http.StatusCreated, response)
}
