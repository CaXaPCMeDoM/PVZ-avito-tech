package app

import (
	"PVZ-avito-tech/config"
	v1 "PVZ-avito-tech/internal/controller/http/v1"
	"PVZ-avito-tech/internal/infrastructure/repo/persistent"
	"PVZ-avito-tech/internal/infrastructure/security/password"
	"PVZ-avito-tech/internal/pkg/auth/jwt"
	"PVZ-avito-tech/internal/pkg/httpserver"
	"PVZ-avito-tech/internal/pkg/logger"
	"PVZ-avito-tech/internal/pkg/postgres"
	"PVZ-avito-tech/internal/usecase/auth"
	"PVZ-avito-tech/internal/usecase/dummy"
	"PVZ-avito-tech/internal/usecase/product"
	"PVZ-avito-tech/internal/usecase/pvz"
	"PVZ-avito-tech/internal/usecase/reception"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	jwtService, err := jwt.NewService([]byte(cfg.Jwt.SecretKey))
	if err != nil {
		l.Fatal("cant create jwt in Run()", err)
	}
	hasher := password.NewBcryptHasher(cfg)

	// repo
	pg, err := postgres.New(cfg.Pg.URL, postgres.MaxPoolSize(cfg.Pg.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	userRepo := persistent.NewUserRepo(pg)
	productRepo := persistent.NewProductRepo(pg)
	receptionRepo := persistent.NewReceptionRepo(pg)
	pvzRepo := persistent.NewPVZRepo(pg)

	// usecase
	userUC := auth.NewUserUsecase(userRepo, hasher)
	dummyUC := dummy.NewDummyAuthUseCase(jwtService)
	pvzUC := pvz.NewPVZUseCase(pvzRepo, receptionRepo, productRepo, l)
	receptionUC := reception.NewUseCase(receptionRepo)
	productUC := product.NewProductUsecase(productRepo)

	// controlerS
	router := v1.NewRouter(
		cfg,
		l,
		userUC,
		dummyUC,
		receptionUC,
		pvzUC,
		productUC,
		jwtService,
	)
	routerMetrics := v1.NewRouterMetrics(
		l,
	)

	// HTTP server
	server := httpserver.New(
		cfg,
		router,
		httpserver.Mode(cfg.HTTP.Mode),
	)

	var serverForMetrics *httpserver.Server
	if cfg.Prometheus.Enabled {
		serverForMetrics = httpserver.New(
			cfg,
			routerMetrics,
			httpserver.Port(cfg.Prometheus.Port),
		)

		serverForMetrics.Start()
	}

	server.Start()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-server.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-server.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = server.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	err = serverForMetrics.Shutdown()

	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
