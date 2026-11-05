package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAuthService struct {
	authenticateFunc func(ctx context.Context, req LinkedInAuthRequest) (*AuthResponse, error)
	refreshFunc      func(ctx context.Context, rawRefreshToken string) (*RefreshTokenResponse, error)
	logoutFunc       func(ctx context.Context, rawRefreshToken string) error
}

func (m *mockAuthService) AuthenticateWithLinkedIn(ctx context.Context, req LinkedInAuthRequest) (*AuthResponse, error) {
	if m.authenticateFunc != nil {
		return m.authenticateFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockAuthService) RefreshToken(ctx context.Context, rawRefreshToken string) (*RefreshTokenResponse, error) {
	if m.refreshFunc != nil {
		return m.refreshFunc(ctx, rawRefreshToken)
	}
	return nil, nil
}

func (m *mockAuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	if m.logoutFunc != nil {
		return m.logoutFunc(ctx, rawRefreshToken)
	}
	return nil
}

func TestHandler_HandleConfig(t *testing.T) {
	handler := NewHandler(nil, "linkedin-client-xyz", "org.nith.alumni://callback")
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/config", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var cfg AuthConfigResponse
	err := json.Unmarshal(rec.Body.Bytes(), &cfg)
	require.NoError(t, err)
	assert.Equal(t, "linkedin-client-xyz", cfg.ClientID)
	assert.Equal(t, "org.nith.alumni://callback", cfg.RedirectURI)
}

func TestHandler_HandleLinkedInCallback(t *testing.T) {
	handler := NewHandler(nil, "client-id", "http://localhost:8080/v1/auth/linkedin/callback")
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/linkedin/callback?code=sample-oauth-code", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "org.nith.alumni://callback?code=sample-oauth-code")
	assert.Contains(t, rec.Body.String(), "Authentication Successful!")
}

func TestHandler_HandleLinkedInAuth(t *testing.T) {
	expectedID := uuid.New()
	mockSvc := &mockAuthService{
		authenticateFunc: func(ctx context.Context, req LinkedInAuthRequest) (*AuthResponse, error) {
			if req.Code == "valid-code" {
				return &AuthResponse{
					AccessToken:  "access-token-123",
					RefreshToken: "refresh-token-123",
					IsNewUser:    true,
					User: AuthUser{
						ID:              expectedID,
						Email:           "new@nith.ac.in",
						Role:            RoleAlumnus,
						Status:          StatusPending,
						ProfileComplete: false,
					},
				}, nil
			}
			return nil, ErrInvalidLinkedInCode
		},
	}

	handler := NewHandler(mockSvc, "test-client-id", "org.nith.alumni://callback")
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	// Valid payload
	body, _ := json.Marshal(LinkedInAuthRequest{
		Code:        "valid-code",
		RedirectURI: "org.nith.alumni://callback",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/linkedin", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp AuthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "access-token-123", resp.AccessToken)
	assert.True(t, resp.IsNewUser)
	assert.Equal(t, expectedID, resp.User.ID)

	// Invalid code
	badBody, _ := json.Marshal(LinkedInAuthRequest{
		Code:        "invalid-code",
		RedirectURI: "org.nith.alumni://callback",
	})
	req = httptest.NewRequest(http.MethodPost, "/v1/auth/linkedin", bytes.NewReader(badBody))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandler_HandleRefreshToken(t *testing.T) {
	mockSvc := &mockAuthService{
		refreshFunc: func(ctx context.Context, rawRefreshToken string) (*RefreshTokenResponse, error) {
			if rawRefreshToken == "valid-refresh-token" {
				return &RefreshTokenResponse{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				}, nil
			}
			return nil, ErrInvalidToken
		},
	}

	handler := NewHandler(mockSvc, "test-client-id", "org.nith.alumni://callback")
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	// Successful refresh
	body, _ := json.Marshal(RefreshTokenRequest{RefreshToken: "valid-refresh-token"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp RefreshTokenResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", resp.AccessToken)
	assert.Equal(t, "new-refresh-token", resp.RefreshToken)

	// Invalid refresh token
	badBody, _ := json.Marshal(RefreshTokenRequest{RefreshToken: "bad-refresh-token"})
	req = httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", bytes.NewReader(badBody))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandler_HandleLogout(t *testing.T) {
	mockSvc := &mockAuthService{
		logoutFunc: func(ctx context.Context, rawRefreshToken string) error {
			return nil
		},
	}

	handler := NewHandler(mockSvc, "test-client-id", "org.nith.alumni://callback")
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	body, _ := json.Marshal(LogoutRequest{RefreshToken: "token-to-revoke"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
