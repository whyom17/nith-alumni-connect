package server_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nith-alumni-connect/backend/internal/auth"
	"github.com/nith-alumni-connect/backend/internal/config"
	"github.com/nith-alumni-connect/backend/internal/platform/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_HealthAndReadiness(t *testing.T) {
	discardLogger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	cfg := &config.Config{
		Port:     "8080",
		Env:      "test",
		LogLevel: "info",
	}

	srv := server.New(server.Params{
		Config:   cfg,
		Logger:   discardLogger,
		DBPinger: nil,
	})

	router := srv.Router()

	t.Run("GET /healthz returns ok with request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Header().Get("X-Request-Id"))

		var body map[string]string
		err := json.NewDecoder(rec.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "ok", body["status"])
	})

	t.Run("GET /readyz returns ready with request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Header().Get("X-Request-Id"))

		var body map[string]string
		err := json.NewDecoder(rec.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "ready", body["status"])
	})

	t.Run("GET /v1/auth/config returns client config", func(t *testing.T) {
		srvWithAuth := server.New(server.Params{
			Config:      cfg,
			Logger:      discardLogger,
			DBPinger:    nil,
			AuthHandler: auth.NewHandler(nil, "sample-client-id", "org.nith.alumni://callback"),
		})

		req := httptest.NewRequest(http.MethodGet, "/v1/auth/config", nil)
		rec := httptest.NewRecorder()
		srvWithAuth.Router().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var body map[string]string
		err := json.NewDecoder(rec.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "sample-client-id", body["client_id"])
		assert.Equal(t, "org.nith.alumni://callback", body["redirect_uri"])
	})
}
