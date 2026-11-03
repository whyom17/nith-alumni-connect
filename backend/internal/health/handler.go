package health

import (
	"context"
	"net/http"
	"time"

	"github.com/nith-alumni-connect/backend/internal/platform/httputil"
)

// Pinger defines an interface for checking connectivity to dependent services.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Handler handles health check endpoints.
type Handler struct {
	db Pinger
}

// NewHandler returns a new health check Handler.
func NewHandler(db Pinger) *Handler {
	return &Handler{db: db}
}

// Healthz is a liveness probe returning 200 OK when the process is up.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// Readyz is a readiness probe checking if downstream dependencies (database) are reachable.
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.db != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := h.db.Ping(ctx); err != nil {
			httputil.WriteError(
				w,
				http.StatusServiceUnavailable,
				"service_unavailable",
				"database is not ready",
				map[string]string{"database": err.Error()},
			)
			return
		}
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}
