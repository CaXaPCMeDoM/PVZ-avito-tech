package pvz

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	er "PVZ-avito-tech/internal/controller/http/errors"
	"PVZ-avito-tech/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Routes) GetPVZList(c *gin.Context) {
	var filter dto.ReceptionFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn(er.ErrInvalidRequestBody)
		dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidRequestBody)
		return
	}

	filter.Apply(
		dto.WithPaginationDefaults(),
	)
	pvzList, err := h.pvzUC.GetPVZWithReceptions(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error(entity.ErrGetPVZList, err)
		dto.ErrorResponse(c, http.StatusBadRequest, entity.ErrGetPVZList.Error())
		return
	}

	c.JSON(http.StatusOK, *pvzList)
}
