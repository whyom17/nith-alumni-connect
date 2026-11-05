package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/nith-alumni-connect/backend/internal/platform/httputil"
)

// Handler handles HTTP requests for auth endpoints.
type Handler struct {
	service     Service
	clientID    string
	redirectURI string
}

// NewHandler creates a new auth HTTP handler.
func NewHandler(service Service, clientID, redirectURI string) *Handler {
	return &Handler{
		service:     service,
		clientID:    clientID,
		redirectURI: redirectURI,
	}
}

// RegisterRoutes registers the auth endpoints on the provided chi router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/v1/auth/config", h.HandleConfig)
	r.Get("/v1/auth/linkedin/callback", h.HandleLinkedInCallback)
	r.Post("/v1/auth/linkedin", h.HandleLinkedInAuth)
	r.Post("/v1/auth/refresh", h.HandleRefreshToken)
	r.Post("/v1/auth/logout", h.HandleLogout)
}

// HandleConfig returns the public client configuration for OAuth initialization.
func (h *Handler) HandleConfig(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, AuthConfigResponse{
		ClientID:    h.clientID,
		RedirectURI: h.redirectURI,
	})
}

// HandleLinkedInCallback receives the redirect from LinkedIn in the browser and forwards to the mobile app scheme.
func (h *Handler) HandleLinkedInCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	errParam := r.URL.Query().Get("error")

	// Android Intent URL (required by Android Chrome to open native apps from HTTP redirects)
	androidIntent := fmt.Sprintf("intent://callback?code=%s&error=%s#Intent;scheme=org.nith.alumni;action=android.intent.action.VIEW;category=android.intent.category.BROWSABLE;package=org.nith.alumni;end;", url.QueryEscape(code), url.QueryEscape(errParam))
	// Standard custom scheme fallback
	customScheme := fmt.Sprintf("org.nith.alumni://callback?code=%s&error=%s", url.QueryEscape(code), url.QueryEscape(errParam))

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Authentication Successful | NITH Alumni Connect</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
      text-align: center;
      padding: 24px 16px;
      background: #FDFBD4;
      color: #3B2712;
      margin: 0;
    }
    .card {
      max-width: 440px;
      margin: 40px auto;
      padding: 32px 24px;
      background: #FFFDEA;
      border-radius: 16px;
      box-shadow: 0 4px 20px rgba(59,39,18,0.1);
      border: 1.5px solid #D47E30;
    }
    h2 { margin-top: 0; color: #3B2712; font-size: 24px; }
    p { color: #5C4028; font-size: 15px; line-height: 1.5; }
    .btn {
      display: inline-block;
      margin-top: 16px;
      padding: 14px 28px;
      background: #8D5A2B;
      color: #FDFBD4 !important;
      text-decoration: none;
      border-radius: 10px;
      font-weight: bold;
      font-size: 16px;
      border: none;
      cursor: pointer;
    }
    .btn:hover { background: #734821; }
    .code-box {
      margin-top: 24px;
      padding: 16px;
      background: #FDFBD4;
      border-radius: 10px;
      border: 1px dashed #D47E30;
    }
    .code-text {
      font-family: monospace;
      font-size: 13px;
      word-break: break-all;
      background: #FFF;
      padding: 8px 12px;
      border-radius: 6px;
      border: 1px solid #D47E30;
      display: block;
      margin: 8px 0;
      user-select: all;
    }
    .copy-btn {
      background: #D47E30;
      color: #3B2712;
      font-weight: bold;
      padding: 8px 16px;
      border-radius: 6px;
      border: none;
      cursor: pointer;
      font-size: 13px;
    }
  </style>
</head>
<body>
  <div class="card">
    <div style="font-size: 44px; margin-bottom: 8px;">🎉</div>
    <h2>Authentication Successful!</h2>
    <p>Returning you to <strong>NITH Alumni Connect</strong>...</p>

    <!-- Action button tries Android intent and custom scheme -->
    <a id="openBtn" class="btn" href="%s">Open App</a>

    <div class="code-box">
      <p style="margin: 0; font-size: 13px;">If the app didn't open automatically (e.g. on Desktop), copy this code and paste it into the app:</p>
      <code id="authCode" class="code-text">%s</code>
      <button class="copy-btn" onclick="copyCode()">📋 Copy Code</button>
      <span id="copyMsg" style="display:none; color: #2e7d32; font-size: 12px; margin-left: 8px; font-weight: bold;">Copied!</span>
    </div>
  </div>

  <script>
    const androidIntent = "%s";
    const customScheme = "%s";

    function copyCode() {
      const code = document.getElementById('authCode').innerText.trim();
      navigator.clipboard.writeText(code).then(() => {
        document.getElementById('copyMsg').style.display = 'inline';
        setTimeout(() => { document.getElementById('copyMsg').style.display = 'none'; }, 2500);
      });
    }

    // Try opening via Android intent or custom scheme
    try {
      window.location.href = androidIntent;
      setTimeout(() => {
        window.location.href = customScheme;
      }, 500);
    } catch (e) {}

    document.getElementById('openBtn').addEventListener('click', function(e) {
      setTimeout(function() {
        window.location.href = customScheme;
      }, 300);
    });
  </script>
</body>
</html>`, androidIntent, code, androidIntent, customScheme)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

// HandleLinkedInAuth exchanges the LinkedIn authorization code for tokens.
func (h *Handler) HandleLinkedInAuth(w http.ResponseWriter, r *http.Request) {
	var req LinkedInAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid_request", "malformed JSON request payload", nil)
		return
	}

	resp, err := h.service.AuthenticateWithLinkedIn(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRequest):
			httputil.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error(), nil)
		case errors.Is(err, ErrInvalidLinkedInCode), errors.Is(err, ErrLinkedInAuthFailed):
			httputil.WriteError(w, http.StatusUnauthorized, "invalid_linkedin_credentials", err.Error(), nil)
		case errors.Is(err, ErrAccountSuspended):
			httputil.WriteError(w, http.StatusForbidden, "account_suspended", "your account has been suspended", nil)
		default:
			httputil.WriteError(w, http.StatusInternalServerError, "internal_error", "an error occurred during authentication", nil)
		}
		return
	}

	httputil.WriteJSON(w, http.StatusOK, resp)
}

// HandleRefreshToken rotates the refresh token and returns new tokens.
func (h *Handler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid_request", "malformed JSON request payload", nil)
		return
	}

	resp, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRequest):
			httputil.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error(), nil)
		case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrExpiredToken), errors.Is(err, ErrTokenRevoked):
			httputil.WriteError(w, http.StatusUnauthorized, "invalid_token", "refresh token is invalid, expired, or revoked", nil)
		case errors.Is(err, ErrAccountSuspended):
			httputil.WriteError(w, http.StatusForbidden, "account_suspended", "your account has been suspended", nil)
		default:
			httputil.WriteError(w, http.StatusInternalServerError, "internal_error", "an error occurred during token refresh", nil)
		}
		return
	}

	httputil.WriteJSON(w, http.StatusOK, resp)
}

// HandleLogout revokes the user's refresh token.
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid_request", "malformed JSON request payload", nil)
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "internal_error", "an error occurred during logout", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
