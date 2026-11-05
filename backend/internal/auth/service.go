package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrAccountSuspended = errors.New("account is suspended")
	ErrTokenRevoked     = errors.New("refresh token has been revoked")
	ErrInvalidRequest   = errors.New("invalid request payload")
)

// Service defines auth business logic.
type Service interface {
	AuthenticateWithLinkedIn(ctx context.Context, req LinkedInAuthRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, rawRefreshToken string) (*RefreshTokenResponse, error)
	Logout(ctx context.Context, rawRefreshToken string) error
}

type authService struct {
	repo           Repository
	tokenManager   TokenManager
	linkedinClient LinkedInClient
}

// NewService creates a new auth service.
func NewService(repo Repository, tokenManager TokenManager, linkedinClient LinkedInClient) Service {
	return &authService{
		repo:           repo,
		tokenManager:   tokenManager,
		linkedinClient: linkedinClient,
	}
}

func (s *authService) AuthenticateWithLinkedIn(ctx context.Context, req LinkedInAuthRequest) (*AuthResponse, error) {
	if req.Code == "" || req.RedirectURI == "" {
		return nil, fmt.Errorf("%w: code and redirect_uri are required", ErrInvalidRequest)
	}

	userInfo, err := s.linkedinClient.ExchangeCode(ctx, req.Code, req.CodeVerifier, req.RedirectURI)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByLinkedInSub(ctx, userInfo.Sub)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}

	isNewUser := false
	if user == nil {
		// New user: provision account with status=pending, role=alumnus
		createdUser, err := s.repo.CreateUser(ctx, userInfo.Sub, userInfo.Email, RoleAlumnus, StatusPending)
		if err != nil {
			return nil, fmt.Errorf("provisioning user: %w", err)
		}
		user = createdUser
		isNewUser = true
	}

	if user.Status == StatusSuspended {
		return nil, ErrAccountSuspended
	}

	userID, err := pgUUIDToUUID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing user id: %w", err)
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(userID, user.Role, user.Status)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	rawRefreshToken, tokenHash, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	expiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	if _, err := s.repo.CreateRefreshToken(ctx, userID, tokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("saving refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		IsNewUser:    isNewUser,
		User: AuthUser{
			ID:              userID,
			Email:           user.Email,
			Role:            user.Role,
			Status:          user.Status,
			ProfileComplete: false, // Phase 1 M2 manages profile completion
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, rawRefreshToken string) (*RefreshTokenResponse, error) {
	if rawRefreshToken == "" {
		return nil, fmt.Errorf("%w: refresh_token is required", ErrInvalidRequest)
	}

	tokenHash := s.tokenManager.HashRefreshToken(rawRefreshToken)
	tokenRecord, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("retrieving refresh token: %w", err)
	}

	if tokenRecord == nil {
		return nil, ErrInvalidToken
	}

	userID, err := pgUUIDToUUID(tokenRecord.UserID)
	if err != nil {
		return nil, fmt.Errorf("parsing user id: %w", err)
	}

	// Token Reuse Detection: If the token was already revoked, someone is reusing an old/compromised token.
	// Invalidate ALL sessions for this user immediately.
	if tokenRecord.RevokedAt.Valid && !tokenRecord.RevokedAt.Time.IsZero() {
		_ = s.repo.RevokeAllUserRefreshTokens(ctx, userID)
		return nil, fmt.Errorf("%w: reuse detected", ErrTokenRevoked)
	}

	// Check if token has expired
	if tokenRecord.ExpiresAt.Valid && tokenRecord.ExpiresAt.Time.Before(time.Now().UTC()) {
		return nil, ErrExpiredToken
	}

	// Revoke this token as part of rotation
	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("revoking old refresh token: %w", err)
	}

	// Fetch current user details from DB to ensure fresh status and role
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidToken
	}
	if user.Status == StatusSuspended {
		return nil, ErrAccountSuspended
	}

	// Generate new access and refresh tokens
	newAccessToken, err := s.tokenManager.GenerateAccessToken(userID, user.Role, user.Status)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	newRawRefreshToken, newTokenHash, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	expiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	if _, err := s.repo.CreateRefreshToken(ctx, userID, newTokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("saving new refresh token: %w", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRawRefreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}
	tokenHash := s.tokenManager.HashRefreshToken(rawRefreshToken)
	return s.repo.RevokeRefreshToken(ctx, tokenHash)
}

func pgUUIDToUUID(p pgtype.UUID) (uuid.UUID, error) {
	if !p.Valid {
		return uuid.Nil, errors.New("invalid uuid")
	}
	return uuid.FromBytes(p.Bytes[:])
}
