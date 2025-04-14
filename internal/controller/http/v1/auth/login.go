package auth

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	er "PVZ-avito-tech/internal/controller/http/errors"
	"PVZ-avito-tech/internal/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Routes) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.ErrorResponse(c, http.StatusUnauthorized, er.ErrInvalidRequestBody)
		return
	}

	email := req.Email
	password := req.Password

	loginResp, err := h.userUC.Login(c.Request.Context(), email, password)

	if err != nil {
		logFields := map[string]interface{}{
			"email":  email,
			"method": "Login",
		}
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			h.logger.Warn(entity.ErrUserNotFound.Error(), logFields)
			dto.ErrorResponse(c, http.StatusUnauthorized, entity.ErrUserNotFound.Error())
		case errors.Is(err, entity.ErrInvalidPassword):
			h.logger.Warn(entity.ErrInvalidPassword.Error(), logFields)
			dto.ErrorResponse(c, http.StatusUnauthorized, entity.ErrInvalidPassword.Error())
		case errors.Is(err, entity.ErrInternal):
			h.logger.Error(entity.ErrInternal.Error(), logFields)
			dto.ErrorResponse(c, http.StatusInternalServerError, entity.ErrInternal.Error())
		default:
			h.logger.Error("unexpected error", logFields)
			dto.ErrorResponse(c, http.StatusInternalServerError, entity.ErrInternal.Error())
		}
		return
	}

	token, errToken := h.dummyUC.GenerateDummyToken(loginResp.Role)

	if errToken != nil {
		h.logger.Error("token generation failed", map[string]interface{}{
			"error":  errToken.Error(),
			"method": "Login",
			"role":   loginResp.Role,
		})
		dto.ErrorResponse(c, http.StatusInternalServerError, entity.ErrInternal.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Token: token,
	})
}
