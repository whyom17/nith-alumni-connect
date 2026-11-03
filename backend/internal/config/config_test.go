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

	cfg := Load()
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Env)
	assert.Contains(t, cfg.DatabaseURL, "postgres://")
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestLoad_CustomEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://user:pass@remote:5432/db")
	t.Setenv("LOG_LEVEL", "debug")

	cfg := Load()
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "production", cfg.Env)
	assert.Equal(t, "postgres://user:pass@remote:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "debug", cfg.LogLevel)
}
