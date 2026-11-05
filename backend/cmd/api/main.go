package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nith-alumni-connect/backend/internal/auth"
	"github.com/nith-alumni-connect/backend/internal/config"
	platformdb "github.com/nith-alumni-connect/backend/internal/platform/db"
	"github.com/nith-alumni-connect/backend/internal/platform/logger"
	"github.com/nith-alumni-connect/backend/internal/platform/server"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, os.Stdout)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := platformdb.NewPool(ctx, cfg.DatabaseURL)
	var authRepo auth.Repository
	if err != nil {
		log.Warn("database connection failed during startup; database-dependent routes will be degraded", "error", err)
	} else {
		defer pool.Close()
		log.Info("connected to PostgreSQL database")
		queries := platformdb.New(pool)
		authRepo = auth.NewRepository(queries)
	}

	tokenManager := auth.NewTokenManager(cfg.JWTSecret)
	linkedinClient := auth.NewLinkedInClient(cfg.LinkedInClientID, cfg.LinkedInClientSecret)
	authService := auth.NewService(authRepo, tokenManager, linkedinClient)
	authHandler := auth.NewHandler(authService, cfg.LinkedInClientID, cfg.LinkedInRedirectURI)

	srv := server.New(server.Params{
		Config:      cfg,
		Logger:      log,
		DBPinger:    pool,
		AuthHandler: authHandler,
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

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("graceful shutdown failed, forcing termination", "error", err)
			_ = srv.Shutdown(context.Background())
			os.Exit(1)
		}
		log.Info("server shutdown complete")
	}
}
