package main

import (
	"fmt"
	"os"
)

func main() {
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		fmt.Println("❌ AUTH_SERVICE_URL is NOT set")
	} else {
		fmt.Println("✅ AUTH_SERVICE_URL:", authServiceURL)
	}
}
