package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthServiceURL loads your external authentication service URL from environment variables
var AuthServiceURL = os.Getenv("AUTH_SERVICE_URL")

// AuthMiddleware validates the JWT token by calling the external authentication service
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token from the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			ctx.Abort()
			return
		}

		// Extract token (Bearer <token>)
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token by making a request to your authentication service
		isValid, err := validateTokenWithAuthService(token)
		if err != nil || !isValid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Proceed if token is valid
		ctx.Next()
	}
}

// validateTokenWithAuthService calls the external authentication service to verify the JWT token
func validateTokenWithAuthService(token string) (bool, error) {
	// Define the request URL (your auth service should have an endpoint for validation)
	url := fmt.Sprintf("%s/validate-token", AuthServiceURL)

	// Create the request with the token in the Authorization header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check if the token is valid (assuming 200 OK means valid)
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, errors.New("invalid token")
}
