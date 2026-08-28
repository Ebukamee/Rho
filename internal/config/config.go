package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	Environment            string
	DatabaseURL            string
	JWTSecret              string
	PaystackSecretKey      string
	CORSOrigins            []string
	GoogleClientID         string
	GoogleClientSecret     string
	GoogleRedirectURL      string
	RateLimitRequests      int
	RateLimitWindowSeconds int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		Port:                   getEnv("PORT", "8080"),
		Environment:            getEnv("ENVIRONMENT", "development"),
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/rho?sslmode=disable"),
		JWTSecret:              getEnv("JWT_SECRET", ""),
		PaystackSecretKey:      getEnv("PAYSTACK_SECRET_KEY", ""),
		CORSOrigins:            splitCSV(getEnv("CORS_ORIGINS", "http://localhost:3000")),
		GoogleClientID:         getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:     getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:      getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
		RateLimitRequests:      atoi(getEnv("RATE_LIMIT_REQUESTS", "60")),
		RateLimitWindowSeconds: atoi(getEnv("RATE_LIMIT_WINDOW_SECONDS", "60")),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func atoi(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil || v <= 0 {
		return 0
	}
	return v
}
