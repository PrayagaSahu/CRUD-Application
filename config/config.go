package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration loaded from environment variables
type Config struct {
	DBDriver    string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	Environment string
}

// Load reads the environment variables and returns a Config struct.
// It fails fast if any required variable is missing.
func Load() *Config {
	// Load .env file (only needed for local development)
	_ = godotenv.Load()

	cfg := &Config{
		DBDriver:    mustGet("DB_DRIVER"),
		DBHost:      mustGet("DB_HOST"),
		DBPort:      mustGet("DB_PORT"),
		DBUser:      mustGet("DB_USER"),
		DBPassword:  mustGet("DB_PASSWORD"),
		DBName:      mustGet("DB_NAME"),
		ServerPort:  mustGet("PORT"),
		Environment: getOrDefault("ENV", "dev"),
	}

	log.Println("✅ Config loaded successfully")
	return cfg
}

// mustGet returns the environment variable or exits if not found
func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("❌ Required environment variable %s is not set", key)
	}
	return val
}

// getOrDefault returns the environment variable or a fallback if not found
func getOrDefault(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("⚠️  %s not set, using default: %s", key, fallback)
		return fallback
	}
	return val
}
