package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WebPort             string
	SelectelUsername    string
	SelectelAccountID   string
	SelectelPassword    string
	SelectelProjectName string
	LogLevel            string // "debug", "info", "warn", "error"
	LogFormat           string // "text", "json"
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		WebPort:             getEnv("WEB_PORT", "8080"),
		SelectelUsername:    getEnv("SELECTEL_USERNAME", ""),
		SelectelAccountID:   getEnv("SELECTEL_ACCOUNT_ID", ""),
		SelectelPassword:    getEnv("SELECTEL_PASSWORD", ""),
		SelectelProjectName: getEnv("SELECTEL_PROJECT_NAME", ""),
		LogLevel:            getEnv("LOG_LEVEL", "INFO"),
		LogFormat:           getEnv("LOG_FORMAT", "TEXT"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
