package products

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	er "PVZ-avito-tech/internal/controller/http/errors"
	"PVZ-avito-tech/internal/controller/http/mapper"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/metrics"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Routes) AddProduct(c *gin.Context) {
	var req dto.PostAddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil || !req.ProductType.IsValidProductType() {
		dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidRequestBody)
		return
	}

	respEntity, err := h.productUC.AddProduct(c.Request.Context(), &req)

	if err != nil {
		h.logger.Warn(err.Error())
		switch {
		case errors.Is(err, entity.ErrNoActiveReception):
			dto.ErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			dto.ErrorResponse(c, http.StatusBadRequest, er.ErrInvalidRequestBody)
		}

		return
	}

	resp := mapper.EntityProductToProductResponse(respEntity)

	metrics.ProductsAdded.Inc()
	c.JSON(http.StatusCreated, resp)
}
