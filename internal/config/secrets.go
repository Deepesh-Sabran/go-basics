package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Secrets holds all sensitive configuration
type Secrets struct {
	JWTSecret    []byte
	DatabaseURL  string
	RedisAddr    string
	Environment  string
}

var AppSecrets *Secrets

// LoadSecrets loads all secrets from environment variables and .env file
func LoadSecrets() {
	// Load .env file if it exists (only in development)
	// In production, secrets should come from environment variables directly
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️  No .env file found, using environment variables only")
		}
	}

	// Validate required secrets
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET environment variable is not set")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("❌ DATABASE_URL environment variable is not set")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default for development
		log.Println("ℹ️  Using default Redis address:", redisAddr)
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
		log.Println("ℹ️  Using development environment")
	}

	AppSecrets = &Secrets{
		JWTSecret:   []byte(jwtSecret),
		DatabaseURL: databaseURL,
		RedisAddr:   redisAddr,
		Environment: environment,
	}

	log.Println("✅ Secrets loaded successfully")
	if environment == "production" {
		log.Println("🔒 Running in PRODUCTION mode")
	}
}

// GetJWTSecret returns the JWT secret as byte slice
func GetJWTSecret() []byte {
	if AppSecrets == nil {
		log.Fatal("❌ Secrets not initialized. Call LoadSecrets() first")
	}
	return AppSecrets.JWTSecret
}

// GetDatabaseURL returns the database connection URL
func GetDatabaseURL() string {
	if AppSecrets == nil {
		log.Fatal("❌ Secrets not initialized. Call LoadSecrets() first")
	}
	return AppSecrets.DatabaseURL
}

// GetRedisAddr returns the Redis address
func GetRedisAddr() string {
	if AppSecrets == nil {
		log.Fatal("❌ Secrets not initialized. Call LoadSecrets() first")
	}
	return AppSecrets.RedisAddr
}

// IsProduction returns true if running in production
func IsProduction() bool {
	if AppSecrets == nil {
		return false
	}
	return AppSecrets.Environment == "production"
}
