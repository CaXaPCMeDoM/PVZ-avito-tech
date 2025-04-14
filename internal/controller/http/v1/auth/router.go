package auth

import (
	"PVZ-avito-tech/internal/pkg/logger"
	"PVZ-avito-tech/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	dummyUC usecase.DummyLogin
	userUC  usecase.Auth
	logger  logger.Interface
}

type tokenResponse struct {
	Token string `json:"token"`
}

func NewAuthRoutes(
	apiV1Group *gin.RouterGroup,
	dummyUC usecase.DummyLogin,
	userUC usecase.Auth,
	logger logger.Interface,
) *Routes {
	au := &Routes{
		dummyUC: dummyUC,
		userUC:  userUC,
		logger:  logger,
	}

	authGroup := apiV1Group.Group("/")
	{
		authGroup.POST("/dummyLogin", au.DummyLogin)
		authGroup.POST("/register", au.Register)
		authGroup.POST("/login", au.Login)
	}

	return au
}
