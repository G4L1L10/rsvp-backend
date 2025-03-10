package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/g4l1l10/rsvp-backend/db"
	"github.com/g4l1l10/rsvp-backend/handlers"
	"github.com/g4l1l10/rsvp-backend/middlewares"
	"github.com/g4l1l10/rsvp-backend/repository"
	"github.com/g4l1l10/rsvp-backend/routes"
	"github.com/g4l1l10/rsvp-backend/service"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// Register API routes
	routes.SetupRoutes(router, guestHandler)

	// Get Cloud Run Port (Cloud Run requires this)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 RSVP Backend is running on port %s...", port)
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("❌ Error starting server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("🛑 Shutting down RSVP backend gracefully...")
}
