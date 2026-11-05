package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nith-alumni-connect/backend/internal/auth"
)

func TestAuthenticate_ValidToken(t *testing.T) {
	tm := auth.NewTokenManager("test-secret-at-least-32-chars-long-key")
	userID := uuid.New()
	tokenStr, err := tm.GenerateAccessToken(userID, auth.RoleAlumnus, auth.StatusApproved)
	require.NoError(t, err)

	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		claims, ok := GetClaims(r.Context())
		assert.True(t, ok)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, auth.RoleAlumnus, claims.Role)
		assert.Equal(t, auth.StatusApproved, claims.Status)
		w.WriteHeader(http.StatusOK)
	})

	handler := Authenticate(tm)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthenticate_MissingAndInvalidToken(t *testing.T) {
	tm := auth.NewTokenManager("test-secret-at-least-32-chars-long-key")
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Authenticate(tm)(nextHandler)

	// Missing header
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Invalid token
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRequireApproved(t *testing.T) {
	tm := auth.NewTokenManager("test-secret-at-least-32-chars-long-key")
	userID := uuid.New()

	approvedToken, err := tm.GenerateAccessToken(userID, auth.RoleAlumnus, auth.StatusApproved)
	require.NoError(t, err)

	pendingToken, err := tm.GenerateAccessToken(userID, auth.RoleAlumnus, auth.StatusPending)
	require.NoError(t, err)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	pipeline := Authenticate(tm)(RequireApproved(nextHandler))

	// Approved user should succeed
	req := httptest.NewRequest(http.MethodGet, "/directory", nil)
	req.Header.Set("Authorization", "Bearer "+approvedToken)
	rec := httptest.NewRecorder()
	pipeline.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Pending user should be forbidden with account_not_approved
	req = httptest.NewRequest(http.MethodGet, "/directory", nil)
	req.Header.Set("Authorization", "Bearer "+pendingToken)
	rec = httptest.NewRecorder()
	pipeline.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "account_not_approved")
}

func TestRequireAdmin(t *testing.T) {
	tm := auth.NewTokenManager("test-secret-at-least-32-chars-long-key")
	userID := uuid.New()

	adminToken, err := tm.GenerateAccessToken(userID, auth.RoleAdmin, auth.StatusApproved)
	require.NoError(t, err)

	alumnusToken, err := tm.GenerateAccessToken(userID, auth.RoleAlumnus, auth.StatusApproved)
	require.NoError(t, err)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	pipeline := Authenticate(tm)(RequireAdmin(nextHandler))

	// Admin succeeds
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec := httptest.NewRecorder()
	pipeline.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Alumnus forbidden
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+alumnusToken)
	rec = httptest.NewRecorder()
	pipeline.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
