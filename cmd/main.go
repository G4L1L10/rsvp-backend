package main

import (
	"fmt"
	"log"

	"github.com/g4l1l10/rsvp-backend/config"
	"github.com/g4l1l10/rsvp-backend/db"
	"github.com/g4l1l10/rsvp-backend/handlers"
	"github.com/g4l1l10/rsvp-backend/middlewares"
	"github.com/g4l1l10/rsvp-backend/repository"
	"github.com/g4l1l10/rsvp-backend/routes"
	"github.com/g4l1l10/rsvp-backend/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Ensure Database URL is available
	if cfg.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL is not set in environment variables")
	}

	// Initialize database connection
	db.InitDB()
	defer db.GetDB().Close()

	// Set up repository, service, and handler
	guestRepo := repository.NewGuestRepository(db.GetDB())
	guestService := service.NewGuestService(guestRepo)
	guestHandler := handlers.NewGuestHandler(guestService)

	// Initialize router with middleware
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware()) // Apply CORS middleware

	// Get server port (fallback to 8080)
	serverPort := cfg.ServerPort
	if serverPort == "" {
		serverPort = "8080"
	}

	// Register API routes
	routes.SetupRoutes(router, guestHandler)

	// Start server
	serverAddr := fmt.Sprintf(":%s", serverPort)
	log.Printf("🚀 Server is running on port %s...", serverPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}
