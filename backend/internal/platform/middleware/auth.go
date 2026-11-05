package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/nith-alumni-connect/backend/internal/auth"
	"github.com/nith-alumni-connect/backend/internal/platform/httputil"
)

type contextKey string

const (
	claimsContextKey contextKey = "auth_claims"
)

// Authenticate verifies the JWT bearer token and injects claims into context.
func Authenticate(tokenManager auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing or malformed Authorization header", nil)
				return
			}

			tokenStr := strings.TrimSpace(authHeader[7:])
			claims, err := tokenManager.ValidateAccessToken(tokenStr)
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid_token", "access token is invalid or expired", nil)
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireApproved ensures that the authenticated user has status "approved".
func RequireApproved(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(claimsContextKey).(*auth.Claims)
		if !ok || claims == nil {
			httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "authentication required", nil)
			return
		}

		if claims.Status != auth.StatusApproved {
			httputil.WriteError(w, http.StatusForbidden, "account_not_approved", "access restricted to approved alumni", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAdmin ensures that the authenticated user has role "admin".
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(claimsContextKey).(*auth.Claims)
		if !ok || claims == nil {
			httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "authentication required", nil)
			return
		}

		if claims.Role != auth.RoleAdmin {
			httputil.WriteError(w, http.StatusForbidden, "forbidden", "admin access required", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetClaims retrieves auth claims from the request context.
func GetClaims(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*auth.Claims)
	return claims, ok
}

// GetUserID retrieves the authenticated user's ID from context.
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	claims, ok := GetClaims(ctx)
	if !ok || claims == nil {
		return uuid.Nil, false
	}
	return claims.UserID, true
}
