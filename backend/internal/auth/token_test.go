package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenManager_GenerateAndValidateAccessToken(t *testing.T) {
	secret := "test-secret-key-at-least-32-chars-long"
	tm := NewTokenManager(secret)
	userID := uuid.New()

	tokenStr, err := tm.GenerateAccessToken(userID, RoleAlumnus, StatusApproved)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := tm.ValidateAccessToken(tokenStr)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, RoleAlumnus, claims.Role)
	assert.Equal(t, StatusApproved, claims.Status)
}

func TestTokenManager_ValidateInvalidToken(t *testing.T) {
	tm := NewTokenManager("secret-1")
	tm2 := NewTokenManager("secret-2")
	userID := uuid.New()

	tokenStr, err := tm.GenerateAccessToken(userID, RoleAlumnus, StatusPending)
	require.NoError(t, err)

	// Validating with wrong secret should fail
	_, err = tm2.ValidateAccessToken(tokenStr)
	assert.ErrorIs(t, err, ErrInvalidToken)

	// Validating malformed token string
	_, err = tm.ValidateAccessToken("invalid.jwt.token")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenManager_ExpiredToken(t *testing.T) {
	manager := &jwtTokenManager{
		secretKey:    []byte("secret-key"),
		accessExpiry: -1 * time.Minute, // already expired
	}

	userID := uuid.New()
	tokenStr, err := manager.GenerateAccessToken(userID, RoleAlumnus, StatusPending)
	require.NoError(t, err)

	_, err = manager.ValidateAccessToken(tokenStr)
	assert.ErrorIs(t, err, ErrExpiredToken)
}

func TestTokenManager_RefreshTokenGenerationAndHash(t *testing.T) {
	tm := NewTokenManager("secret-key")

	rawToken1, hash1, err := tm.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEmpty(t, rawToken1)
	assert.NotEmpty(t, hash1)

	// Hashing the raw token again produces the identical hash
	assert.Equal(t, hash1, tm.HashRefreshToken(rawToken1))

	rawToken2, hash2, err := tm.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEqual(t, rawToken1, rawToken2)
	assert.NotEqual(t, hash1, hash2)
}
