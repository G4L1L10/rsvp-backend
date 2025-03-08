package middlewares

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token by calling the authentication service
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ✅ Fetch AUTH_SERVICE_URL inside the function to ensure it is read at runtime
		authServiceURL := os.Getenv("AUTH_SERVICE_URL")
		if authServiceURL == "" {
			log.Println("❌ AUTH_SERVICE_URL is not set")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "AUTH_SERVICE_URL is not set"})
			ctx.Abort()
			return
		}

		// Get token from the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			ctx.Abort()
			return
		}

		// Extract token (Bearer <token>)
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token with the authentication service
		valid, err := validateTokenWithAuthService(authServiceURL, token)
		if err != nil || !valid {
			log.Printf("❌ Token validation failed: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Proceed if token is valid
		ctx.Next()
	}
}

// validateTokenWithAuthService sends the token to the authentication service for validation
func validateTokenWithAuthService(authServiceURL, token string) (bool, error) {
	// ✅ Ensure the URL is not empty
	if authServiceURL == "" {
		log.Println("❌ AUTH_SERVICE_URL is empty in validateTokenWithAuthService")
		return false, fmt.Errorf("AUTH_SERVICE_URL is empty")
	}

	// Define the validation URL
	url := fmt.Sprintf("%s/auth/validate", authServiceURL)

	// Create the request with the token in the Authorization header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Log the outgoing request
	log.Printf("🔄 Validating token with auth service: %s", url)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error calling auth service: %v", err)
		return false, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Error reading auth service response: %v", err)
		return false, err
	}

	// Log the response
	log.Printf("🔄 Auth service response: %s", string(body))

	// If status is 200, the token is valid
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("invalid token")
}
