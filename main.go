package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boichique/eleven_test_task/internal/config"
	"github.com/boichique/eleven_test_task/internal/server"
	"golang.org/x/exp/slog"
)

const (
	gracefulTimeout = 10 * time.Second
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	srv := server.NewServer(cfg)

	go gracefulShutdown(srv)

	if err := srv.Start(); err != http.ErrServerClosed {
		slog.Error(
			"server stopped",
			"error", err,
		)
		os.Exit(1)
	}
}

func failOnError(err error, message string) {
	if err != nil {
		slog.Error("%s: %s", message, err)
		os.Exit(1)
	}
}

func gracefulShutdown(srv *server.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	<-sigCh
	slog.Info("received interrupt signal. Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error(
			"shutdown server",
			"error", err,
		)
	}
}
