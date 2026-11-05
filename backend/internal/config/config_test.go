package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("ENV")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("LINKEDIN_CLIENT_ID")
	os.Unsetenv("LINKEDIN_CLIENT_SECRET")
	os.Unsetenv("LINKEDIN_REDIRECT_URI")

	cfg := Load()
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Env)
	assert.Contains(t, cfg.DatabaseURL, "postgres://")
	assert.Equal(t, "info", cfg.LogLevel)
	assert.NotEmpty(t, cfg.JWTSecret)
	assert.Equal(t, "http://localhost:8080/v1/auth/linkedin/callback", cfg.LinkedInRedirectURI)
}

func TestLoad_CustomEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://user:pass@remote:5432/db")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("JWT_SECRET", "custom-super-secret-key-1234567890123456")
	t.Setenv("LINKEDIN_CLIENT_ID", "test-client-id")
	t.Setenv("LINKEDIN_CLIENT_SECRET", "test-client-secret")
	t.Setenv("LINKEDIN_REDIRECT_URI", "https://alumni.nith.ac.in/callback")

	cfg := Load()
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "production", cfg.Env)
	assert.Equal(t, "postgres://user:pass@remote:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "custom-super-secret-key-1234567890123456", cfg.JWTSecret)
	assert.Equal(t, "test-client-id", cfg.LinkedInClientID)
	assert.Equal(t, "test-client-secret", cfg.LinkedInClientSecret)
	assert.Equal(t, "https://alumni.nith.ac.in/callback", cfg.LinkedInRedirectURI)
}
