package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nith-alumni-connect/backend/internal/platform/db"
)

type mockLinkedInClient struct {
	userInfo *LinkedInUserInfo
	err      error
}

func (m *mockLinkedInClient) ExchangeCode(ctx context.Context, code, codeVerifier, redirectURI string) (*LinkedInUserInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.userInfo, nil
}

type mockRepository struct {
	users             map[string]*db.User        // key: linkedin_sub
	usersByID         map[uuid.UUID]*db.User     // key: uuid
	tokens            map[string]*db.RefreshToken // key: token_hash
	allTokensRevoked  []uuid.UUID
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		users:     make(map[string]*db.User),
		usersByID: make(map[uuid.UUID]*db.User),
		tokens:    make(map[string]*db.RefreshToken),
	}
}

func (m *mockRepository) GetUserByLinkedInSub(ctx context.Context, sub string) (*db.User, error) {
	return m.users[sub], nil
}

func (m *mockRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	return m.usersByID[id], nil
}

func (m *mockRepository) CreateUser(ctx context.Context, linkedinSub, email, role, status string) (*db.User, error) {
	newID := uuid.New()
	var pgUUID pgtype.UUID
	_ = pgUUID.Scan(newID.String())

	u := &db.User{
		ID:          pgUUID,
		LinkedinSub: linkedinSub,
		Email:       email,
		Role:        role,
		Status:      status,
	}
	m.users[linkedinSub] = u
	m.usersByID[newID] = u
	return u, nil
}

func (m *mockRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (*db.RefreshToken, error) {
	var pgUUID pgtype.UUID
	_ = pgUUID.Scan(userID.String())

	tok := &db.RefreshToken{
		UserID:    pgUUID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	}
	m.tokens[tokenHash] = tok
	return tok, nil
}

func (m *mockRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*db.RefreshToken, error) {
	return m.tokens[tokenHash], nil
}

func (m *mockRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	if tok, ok := m.tokens[tokenHash]; ok {
		tok.RevokedAt = pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	}
	return nil
}

func (m *mockRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	m.allTokensRevoked = append(m.allTokensRevoked, userID)
	for _, tok := range m.tokens {
		id, _ := pgUUIDToUUID(tok.UserID)
		if id == userID {
			tok.RevokedAt = pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
		}
	}
	return nil
}

func TestAuthService_AuthenticateWithLinkedIn_NewUser(t *testing.T) {
	repo := newMockRepository()
	tm := NewTokenManager("test-secret-key-at-least-32-chars-long")
	liClient := &mockLinkedInClient{
		userInfo: &LinkedInUserInfo{
			Sub:   "linkedin-123",
			Email: "test@nith.ac.in",
			Name:  "Test Alumnus",
		},
	}

	service := NewService(repo, tm, liClient)
	resp, err := service.AuthenticateWithLinkedIn(context.Background(), LinkedInAuthRequest{
		Code:        "auth-code",
		RedirectURI: "org.nith.alumni://callback",
	})

	require.NoError(t, err)
	assert.True(t, resp.IsNewUser)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, "test@nith.ac.in", resp.User.Email)
	assert.Equal(t, StatusPending, resp.User.Status)
	assert.Equal(t, RoleAlumnus, resp.User.Role)
}

func TestAuthService_AuthenticateWithLinkedIn_ExistingUser(t *testing.T) {
	repo := newMockRepository()
	tm := NewTokenManager("test-secret-key-at-least-32-chars-long")

	// Pre-seed an existing approved user
	existingUser, err := repo.CreateUser(context.Background(), "linkedin-123", "existing@nith.ac.in", RoleAlumnus, StatusApproved)
	require.NoError(t, err)

	liClient := &mockLinkedInClient{
		userInfo: &LinkedInUserInfo{
			Sub:   "linkedin-123",
			Email: "existing@nith.ac.in",
		},
	}

	service := NewService(repo, tm, liClient)
	resp, err := service.AuthenticateWithLinkedIn(context.Background(), LinkedInAuthRequest{
		Code:        "auth-code",
		RedirectURI: "org.nith.alumni://callback",
	})

	require.NoError(t, err)
	assert.False(t, resp.IsNewUser)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	existingID, _ := pgUUIDToUUID(existingUser.ID)
	assert.Equal(t, existingID, resp.User.ID)
	assert.Equal(t, StatusApproved, resp.User.Status)
}

func TestAuthService_AuthenticateWithLinkedIn_SuspendedUser(t *testing.T) {
	repo := newMockRepository()
	tm := NewTokenManager("test-secret-key-at-least-32-chars-long")

	_, err := repo.CreateUser(context.Background(), "linkedin-suspended", "bad@nith.ac.in", RoleAlumnus, StatusSuspended)
	require.NoError(t, err)

	liClient := &mockLinkedInClient{
		userInfo: &LinkedInUserInfo{
			Sub:   "linkedin-suspended",
			Email: "bad@nith.ac.in",
		},
	}

	service := NewService(repo, tm, liClient)
	_, err = service.AuthenticateWithLinkedIn(context.Background(), LinkedInAuthRequest{
		Code:        "auth-code",
		RedirectURI: "org.nith.alumni://callback",
	})

	assert.ErrorIs(t, err, ErrAccountSuspended)
}

func TestAuthService_RefreshToken_Rotation(t *testing.T) {
	repo := newMockRepository()
	tm := NewTokenManager("test-secret-key-at-least-32-chars-long")
	service := NewService(repo, tm, nil)

	user, err := repo.CreateUser(context.Background(), "sub-1", "user@nith.ac.in", RoleAlumnus, StatusApproved)
	require.NoError(t, err)
	userID, _ := pgUUIDToUUID(user.ID)

	rawToken, tokenHash, err := tm.GenerateRefreshToken()
	require.NoError(t, err)
	_, err = repo.CreateRefreshToken(context.Background(), userID, tokenHash, time.Now().UTC().Add(24*time.Hour))
	require.NoError(t, err)

	// Perform refresh
	refreshResp, err := service.RefreshToken(context.Background(), rawToken)
	require.NoError(t, err)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.NotEmpty(t, refreshResp.RefreshToken)
	assert.NotEqual(t, rawToken, refreshResp.RefreshToken)

	// Old token should now be revoked
	oldTokenRecord, err := repo.GetRefreshTokenByHash(context.Background(), tokenHash)
	require.NoError(t, err)
	assert.True(t, oldTokenRecord.RevokedAt.Valid)
}

func TestAuthService_RefreshToken_ReuseDetection(t *testing.T) {
	repo := newMockRepository()
	tm := NewTokenManager("test-secret-key-at-least-32-chars-long")
	service := NewService(repo, tm, nil)

	user, err := repo.CreateUser(context.Background(), "sub-reuse", "reuse@nith.ac.in", RoleAlumnus, StatusApproved)
	require.NoError(t, err)
	userID, _ := pgUUIDToUUID(user.ID)

	rawToken, tokenHash, err := tm.GenerateRefreshToken()
	require.NoError(t, err)
	_, err = repo.CreateRefreshToken(context.Background(), userID, tokenHash, time.Now().UTC().Add(24*time.Hour))
	require.NoError(t, err)

	// First refresh: succeeds and revokes rawToken
	_, err = service.RefreshToken(context.Background(), rawToken)
	require.NoError(t, err)

	// Second refresh using the SAME (now revoked) token: REUSE DETECTED!
	_, err = service.RefreshToken(context.Background(), rawToken)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrTokenRevoked))

	// All tokens for this user must have been revoked
	assert.Contains(t, repo.allTokensRevoked, userID)
}

func TestAuthService_Logout(t *testing.T) {
	repo := newMockRepository()
	tm := NewTokenManager("test-secret-key-at-least-32-chars-long")
	service := NewService(repo, tm, nil)

	userID := uuid.New()
	rawToken, tokenHash, err := tm.GenerateRefreshToken()
	require.NoError(t, err)
	_, err = repo.CreateRefreshToken(context.Background(), userID, tokenHash, time.Now().UTC().Add(24*time.Hour))
	require.NoError(t, err)

	err = service.Logout(context.Background(), rawToken)
	require.NoError(t, err)

	tok, err := repo.GetRefreshTokenByHash(context.Background(), tokenHash)
	require.NoError(t, err)
	assert.True(t, tok.RevokedAt.Valid)
}
