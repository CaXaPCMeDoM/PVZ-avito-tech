package pvz

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	er "PVZ-avito-tech/internal/controller/http/errors"
	"PVZ-avito-tech/internal/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Routes) CloseReception(c *gin.Context) {
	pvzIdStr := c.Param("pvzId")
	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidParam)
		return
	}

	response, err := h.receptionUC.CloseReception(c.Request.Context(), pvzId)

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoActiveReception):
			h.logger.Warn(err.Error())
			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			h.logger.Error(err.Error())
			dto.ErrorResponse(c, http.StatusInternalServerError, entity.ErrInternal.Error())
		}
		return
	}

	c.JSON(http.StatusOK, response)
}
