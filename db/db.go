package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	// Read DATABASE_URL from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL is not set in environment variables")
	}

	// Open connection
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Error connecting to the database: %v", err)
	}

	// Verify connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	log.Println("✅ Connected to CockroachDB successfully!")
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
}
