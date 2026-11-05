package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Config holds runtime configuration populated from environment variables.
type Config struct {
	Env                  string
	Port                 string
	DatabaseURL          string
	LogLevel             string
	JWTSecret            string
	LinkedInClientID     string
	LinkedInClientSecret string
	LinkedInRedirectURI  string
}

// Load loads the configuration from environment variables (loading .env if present) with sensible defaults.
func Load() *Config {
	loadDotEnv(".env", "../.env", "backend/.env")

	return &Config{
		Env:                  getEnv("ENV", "development"),
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/alumni_connect?sslmode=disable"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		JWTSecret:            getEnv("JWT_SECRET", "dev-secret-key-at-least-32-characters-long-1234567890"),
		LinkedInClientID:     getEnv("LINKEDIN_CLIENT_ID", ""),
		LinkedInClientSecret: getEnv("LINKEDIN_CLIENT_SECRET", ""),
		LinkedInRedirectURI:  getEnv("LINKEDIN_REDIRECT_URI", "http://localhost:8080/v1/auth/linkedin/callback"),
	}
}

func loadDotEnv(paths ...string) {
	for _, p := range paths {
		file, err := os.Open(p)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				// Strip quotes if present
				if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
					val = val[1 : len(val)-1]
				}
				// Set in environment if not already explicitly set
				if _, exists := os.LookupEnv(key); !exists {
					_ = os.Setenv(key, val)
				}
			}
		}
		return // Stop after successfully loading the first matching .env file
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
