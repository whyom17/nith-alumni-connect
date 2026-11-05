package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// User role constants.
const (
	RoleAlumnus = "alumnus"
	RoleAdmin   = "admin"
)

// User status constants.
const (
	StatusPending   = "pending"
	StatusApproved  = "approved"
	StatusRejected  = "rejected"
	StatusSuspended = "suspended"
)

// AuthConfigResponse provides public client configuration for starting OAuth flows.
type AuthConfigResponse struct {
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
}

// LinkedInAuthRequest is the payload sent from the mobile client.
type LinkedInAuthRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier,omitempty"`
	RedirectURI  string `json:"redirect_uri"`
}

// AuthUser represents the authenticated user payload in responses.
type AuthUser struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	Status          string    `json:"status"`
	ProfileComplete bool      `json:"profile_complete"`
}

// AuthResponse is returned on successful LinkedIn login/signup.
type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	IsNewUser    bool     `json:"is_new_user"`
	User         AuthUser `json:"user"`
}

// RefreshTokenRequest is the payload for token rotation.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse contains the newly rotated tokens.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest is the payload to terminate a session.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Claims represents the JWT claims issued by our application.
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	Status string    `json:"status"`
}

// LinkedInUserInfo represents user profile data obtained from LinkedIn OIDC.
type LinkedInUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// TokenSession represents a stored refresh token entity.
type TokenSession struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	RevokedAt *time.Time
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
