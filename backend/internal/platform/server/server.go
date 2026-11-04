package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/nith-alumni-connect/backend/internal/config"
	"github.com/nith-alumni-connect/backend/internal/health"
	appmiddleware "github.com/nith-alumni-connect/backend/internal/platform/middleware"
)

// Server wraps http.Server with chi router and logging.
type Server struct {
	cfg        *config.Config
	log        *slog.Logger
	router     *chi.Mux
	httpServer *http.Server
}

// Params holds dependencies for creating a Server.
type Params struct {
	Config   *config.Config
	Logger   *slog.Logger
	DBPinger health.Pinger
}

// New constructs a new Server and configures routes and middleware.
func New(p Params) *Server {
	r := chi.NewRouter()

	// Base middleware stack
	r.Use(appmiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(appmiddleware.Recoverer(p.Logger))
	r.Use(appmiddleware.RequestLogger(p.Logger))
	r.Use(appmiddleware.CORS([]string{"*"}))

	// Health and readiness endpoints
	healthHandler := health.NewHandler(p.DBPinger)
	r.Get("/healthz", healthHandler.Healthz)
	r.Get("/readyz", healthHandler.Readyz)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", p.Config.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		cfg:        p.Config,
		log:        p.Logger,
		router:     r,
		httpServer: httpServer,
	}
}

// Router returns the underlying Chi router (useful for mounting feature routes and in tests).
func (s *Server) Router() *chi.Mux {
	return s.router
}

// Start begins listening for requests.
func (s *Server) Start() error {
	s.log.Info("starting HTTP server",
		slog.String("port", s.cfg.Port),
		slog.String("env", s.cfg.Env),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server failure: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server with a timeout.
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}
