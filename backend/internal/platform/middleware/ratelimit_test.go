package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_LimitsAfterMaxRequests(t *testing.T) {
	limiter := NewRateLimiter(3, 1*time.Minute)
	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	makeReq := func() int {
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec.Code
	}

	// First 3 requests should pass (tokens: 2, 1, 0)
	assert.Equal(t, http.StatusOK, makeReq())
	assert.Equal(t, http.StatusOK, makeReq())
	assert.Equal(t, http.StatusOK, makeReq())

	// 4th request should be rate-limited (429)
	assert.Equal(t, http.StatusTooManyRequests, makeReq())
}
