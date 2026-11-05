package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// TokenManager handles issuing and verifying JWT access tokens and hashing refresh tokens.
type TokenManager interface {
	GenerateAccessToken(userID uuid.UUID, role, status string) (string, error)
	ValidateAccessToken(tokenStr string) (*Claims, error)
	GenerateRefreshToken() (string, string, error) // returns (rawToken, tokenHash, error)
	HashRefreshToken(tokenStr string) string
}

type jwtTokenManager struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewTokenManager creates a TokenManager with default 15-minute access token lifespan.
func NewTokenManager(secretKey string) TokenManager {
	return &jwtTokenManager{
		secretKey:     []byte(secretKey),
		accessExpiry:  15 * time.Minute,
		refreshExpiry: 30 * 24 * time.Hour,
	}
}

func (m *jwtTokenManager) GenerateAccessToken(userID uuid.UUID, role, status string) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    "nith-alumni-connect",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpiry)),
		},
		UserID: userID,
		Role:   role,
		Status: status,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

func (m *jwtTokenManager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (m *jwtTokenManager) GenerateRefreshToken() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("generating random bytes: %w", err)
	}

	rawToken := base64.RawURLEncoding.EncodeToString(bytes)
	tokenHash := m.HashRefreshToken(rawToken)
	return rawToken, tokenHash, nil
}

func (m *jwtTokenManager) HashRefreshToken(tokenStr string) string {
	hash := sha256.Sum256([]byte(tokenStr))
	return hex.EncodeToString(hash[:])
}
