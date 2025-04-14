package pvz

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	er "PVZ-avito-tech/internal/controller/http/errors"
	"PVZ-avito-tech/internal/controller/http/mapper"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Routes) CreatePVZ(c *gin.Context) {
	var req dto.CreatePVZRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(er.ErrInvalidRequestBody)
		dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidRequestBody)
		return
	}

	pvzEntity := mapper.DtoPVZToEntityPVZ(req)

	if !pvzEntity.City.IsValidCity() {
		h.logger.Warn(er.ErrInvalidRequestBody)
		dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidRequestBody)
		return
	}

	pvzResp, err := h.pvzUC.CreatePVZ(c.Request.Context(), pvzEntity)

	if err != nil {
		h.logger.Warn(err.Error())
		dto.ErrorResponse(c, http.StatusBadRequest, entity.ErrCreatePVZ.Error())
		return
	}

	metrics.PVZCreated.Inc()
	c.JSON(http.StatusCreated, pvzResp)
}
