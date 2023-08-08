package server

import (
	"context"
	"fmt"

	"github.com/boichique/eleven_test_task/internal/config"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e   *echo.Echo
	cfg *config.Config
	// ... other fields as needed ...
}

func NewServer(cfg *config.Config) *Server {
	e := echo.New()
	handler := NewHandler()

	api := e.Group("/api")
	api.GET("/items/limits", handler.GetLimits)
	api.POST("/items/process", handler.ProcessBatch)

	return &Server{
		e:   e,
		cfg: cfg,
		// ... initialize other fields as needed ...
	}
}

func (s *Server) Start() error {
	port := s.cfg.Port

	return s.e.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
