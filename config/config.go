package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds environment variables
type Config struct {
	ServerPort     string
	AuthServiceURL string
	DatabaseURL    string
	SMTPUser       string
	SMTPPassword   string
}

// LoadConfig loads environment variables from .env file
func LoadConfig() *Config {
	err := godotenv.Load(".env") // Force loading .env explicitly
	if err != nil {
		log.Println("⚠️ Warning: No .env file found, using system environment variables.")
	}

	// Debugging: Log the values being loaded
	log.Printf("🔄 ENV: SERVER_PORT=%s, AUTH_SERVICE_URL=%s", os.Getenv("SERVER_PORT"), os.Getenv("AUTH_SERVICE_URL"))

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8081"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", ""),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
