package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct to hold application configurations
type Config struct {
	ServerPort  string
	DatabaseURL string
}

// LoadConfig reads environment variables and loads them into a Config struct
func LoadConfig() *Config {
	// Load .env file (optional, but useful in development)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: No .env file found, using system environment variables.")
	}

	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
