package httpserver

import (
	"PVZ-avito-tech/config"
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(cfg *config.Config, router *gin.Engine, opts ...Option) *Server {
	httpServer := &http.Server{
		Addr:           net.JoinHostPort("", cfg.HTTP.Port),
		Handler:        router,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		IdleTimeout:    cfg.HTTP.IdleTimeout,
		MaxHeaderBytes: cfg.HTTP.MaxHeaderBytes,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: cfg.HTTP.ShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	gin.SetMode(cfg.HTTP.Mode)

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
