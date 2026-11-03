package config

import (
	"os"
	"strconv"
)

// Config holds runtime configuration populated from environment variables.
type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	LogLevel    string
}

// Load loads the configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Env:         getEnv("ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/alumni_connect?sslmode=disable"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}
