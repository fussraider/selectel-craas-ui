package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	WebPort             string
	SelectelUsername    string
	SelectelAccountID   string
	SelectelPassword    string
	SelectelProjectName string
	SelectelAuthURL     string
	SelectelProjURL     string
	SelectelCraasURL    string
	LogLevel            string // "debug", "info", "warn", "error"
	LogFormat           string // "text", "json"

	// Feature Flags
	EnableDeleteRegistry   bool
	EnableDeleteRepository bool
	EnableDeleteImage      bool
	EnableMissingTagsCheck bool

	ProtectedTags []string

	// Authentication
	AuthEnabled  bool
	AuthLogin    string
	AuthPassword string
	JWTSecret    string

	// CORS
	CORSAllowedOrigin string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			log.Fatalf("Failed to generate random JWT secret: %v", err)
		}
		jwtSecret = hex.EncodeToString(b)
		log.Println("JWT_SECRET not provided, generated a random secure secret for this session")
	}

	return &Config{
		WebPort:             getEnv("WEB_PORT", "8080"),
		SelectelUsername:    getEnv("SELECTEL_USERNAME", ""),
		SelectelAccountID:   getEnv("SELECTEL_ACCOUNT_ID", ""),
		SelectelPassword:    getEnv("SELECTEL_PASSWORD", ""),
		SelectelProjectName: getEnv("SELECTEL_PROJECT_NAME", ""),
		SelectelAuthURL:     getEnv("SELECTEL_AUTH_URL", "https://cloud.api.selcloud.ru/identity/v3/auth/tokens"),
		SelectelProjURL:     getEnv("SELECTEL_PROJ_URL", "https://cloud.api.selcloud.ru/identity/v3/auth/projects"),
		SelectelCraasURL:    getEnv("SELECTEL_CRAAS_URL", "https://cr.selcloud.ru/api/v1"),
		LogLevel:            getEnv("LOG_LEVEL", "INFO"),
		LogFormat:           getEnv("LOG_FORMAT", "TEXT"),

		EnableDeleteRegistry:   getEnvBool("ENABLE_DELETE_REGISTRY", false),
		EnableDeleteRepository: getEnvBool("ENABLE_DELETE_REPOSITORY", false),
		EnableDeleteImage:      getEnvBool("ENABLE_DELETE_IMAGE", false),
		EnableMissingTagsCheck: getEnvBool("ENABLE_MISSING_TAGS_CHECK", false),

		ProtectedTags: getEnvSlice("PROTECTED_TAGS", nil),

		AuthEnabled:  getEnvBool("AUTH_ENABLED", false),
		AuthLogin:    getEnv("AUTH_LOGIN", ""),
		AuthPassword: getEnv("AUTH_PASSWORD", ""),
		JWTSecret:    jwtSecret,

		CORSAllowedOrigin: getEnv("CORS_ALLOWED_ORIGIN", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		parts := strings.Split(value, ",")
		var trimmed []string
		for _, p := range parts {
			if t := strings.TrimSpace(p); t != "" {
				trimmed = append(trimmed, t)
			}
		}
		return trimmed
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
		log.Printf("Warning: Invalid boolean value for env %s: %s. Using fallback %v", key, value, fallback)
	}
	return fallback
}
