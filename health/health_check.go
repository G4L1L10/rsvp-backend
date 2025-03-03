package health

import (
	"net/http"
	"os"

	"github.com/g4l1l10/rsvp-backend/db"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler returns system health status
func HealthCheckHandler(ctx *gin.Context) {
	status := checkSystemHealth()
	ctx.JSON(http.StatusOK, status)
}

// checkSystemHealth verifies database and authentication service health
func checkSystemHealth() gin.H {
	// Check database connection
	dbStatus := "OK"
	if err := db.GetDB().Ping(); err != nil {
		dbStatus = "ERROR: " + err.Error()
	}

	// Check authentication service
	authStatus := "OK"
	authServiceURL := os.Getenv("AUTH_SERVICE_URL") + "/status"
	resp, err := http.Get(authServiceURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		authStatus = "ERROR: Auth service unavailable"
	}

	return gin.H{
		"server":                 "OK",
		"database":               dbStatus,
		"authentication_service": authStatus,
	}
}
