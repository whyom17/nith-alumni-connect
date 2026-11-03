package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/nith-alumni-connect/backend/internal/platform/httputil"
)

// RequestID sets a unique request ID in context and adds X-Request-Id to the response headers.
func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		if reqID != "" {
			w.Header().Set("X-Request-Id", reqID)
		}
		next.ServeHTTP(w, r)
	}))
}

// RequestLogger logs each request's method, path, status, latency, and request ID.
func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				latency := time.Since(start)
				reqID := middleware.GetReqID(r.Context())

				log.InfoContext(r.Context(), "http request",
					slog.String("request_id", reqID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.Duration("duration", latency),
					slog.String("remote_addr", r.RemoteAddr),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

// Recoverer catches panics, logs them with stack traces, and returns a 500 JSON error.
func Recoverer(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						panic(rvr)
					}

					reqID := middleware.GetReqID(r.Context())
					log.ErrorContext(r.Context(), "panic recovered",
						slog.String("request_id", reqID),
						slog.Any("error", rvr),
						slog.String("stack", string(debug.Stack())),
					)

					httputil.WriteError(
						w,
						http.StatusInternalServerError,
						"internal_error",
						"An unexpected server error occurred",
						nil,
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
