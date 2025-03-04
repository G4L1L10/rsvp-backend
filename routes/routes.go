package routes

import (
	"github.com/g4l1l10/rsvp-backend/handlers"
	"github.com/g4l1l10/rsvp-backend/middlewares"

	"github.com/g4l1l10/rsvp-backend/health"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers API endpoints
func SetupRoutes(router *gin.Engine, guestHandler *handlers.GuestHandler) {
	// Health check route
	router.GET("/status", health.HealthCheckHandler)

	// Public RSVP Routes (Guests can only submit their RSVP)
	rsvpRoutes := router.Group("/rsvp")
	{
		rsvpRoutes.POST("/", guestHandler.SubmitRSVP)
		rsvpRoutes.GET("/:token", guestHandler.GetGuestByToken)
	}

	// Admin Guest Management (Protected)
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middlewares.AuthMiddleware()) // Require JWT authentication
	{
		adminRoutes.GET("/guests", guestHandler.GetAllGuests)
		adminRoutes.GET("/guests/:id", guestHandler.GetGuestByID)
		adminRoutes.GET("/guests/email/:email", guestHandler.GetGuestByEmail)
		adminRoutes.GET("/guests/rsvp/:token", guestHandler.GetGuestByToken)
		adminRoutes.PUT("/guests/:id", guestHandler.UpdateGuest)
		adminRoutes.DELETE("/guests/:id", guestHandler.DeleteGuest)
	}
}
