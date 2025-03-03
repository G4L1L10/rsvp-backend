package routes

import (
	"github.com/g4l1l10/rsvp-backend/handlers"
	"github.com/g4l1l10/rsvp-backend/health"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all API endpoints
func SetupRoutes(router *gin.Engine, guestHandler *handlers.GuestHandler) {
	// Register health check endpoint
	router.GET("/status", health.HealthCheckHandler)

	guestRoutes := router.Group("/guests")
	{
		guestRoutes.POST("/", guestHandler.AddGuest)
		guestRoutes.GET("/", guestHandler.GetAllGuests)
		guestRoutes.GET("/:id", guestHandler.GetGuestByID)
		guestRoutes.GET("/email/:email", guestHandler.GetGuestByEmail)
		guestRoutes.GET("/rsvp/:token", guestHandler.GetGuestByToken)
		guestRoutes.PUT("/:id", guestHandler.UpdateGuest)
		guestRoutes.DELETE("/:id", guestHandler.DeleteGuest)
	}
}
