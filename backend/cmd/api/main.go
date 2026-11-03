package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nith-alumni-connect/backend/internal/config"
	"github.com/nith-alumni-connect/backend/internal/platform/logger"
	"github.com/nith-alumni-connect/backend/internal/platform/server"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, os.Stdout)

	srv := server.New(server.Params{
		Config:   cfg,
		Logger:   log,
		DBPinger: nil, // Will be wired to db.Pool in the database step
	})

	// Start server in background goroutine
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- srv.Start()
	}()

	// Listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error encountered", "error", err)
			os.Exit(1)
		}
	case sig := <-shutdown:
		log.Info("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("graceful shutdown failed, forcing termination", "error", err)
			_ = srv.Shutdown(context.Background())
			os.Exit(1)
		}
		log.Info("server shutdown complete")
	}
}
