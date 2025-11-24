package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	// Port on which the server will listen
	serverPort = ":8080"
	
	// Rate limiting configuration
	maxRequestsPerMinute = 5
	rateLimitWindow      = 1 * time.Minute
)

func main() {
	// Initialize the rate limiter with 5 requests per minute
	// This creates a new RateLimiter instance and starts the background reset goroutine
	rateLimiter := NewRateLimiter(maxRequestsPerMinute, rateLimitWindow)
	
	// Create the middleware function that will wrap our handlers
	// This middleware will check rate limits before allowing requests through
	rateLimitMiddleware := RateLimitMiddleware(rateLimiter)
	
	// Register the /api/data endpoint with rate limiting middleware
	// The middleware runs first, then if allowed, the handler executes
	http.HandleFunc("/api/data", rateLimitMiddleware(handleDataRequest))
	
	// Start the HTTP server
	// This blocks until the server is stopped
	fmt.Printf("Rate Limiter API Server starting on port %s\n", serverPort)
	fmt.Printf("Rate limit: %d requests per %v\n", maxRequestsPerMinute, rateLimitWindow)
	fmt.Println("Try: curl http://localhost:8080/api/data")
	
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleDataRequest is the handler for the /api/data endpoint
// This function only executes if the rate limit middleware allows the request
func handleDataRequest(w http.ResponseWriter, r *http.Request) {
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	
	// Only allow GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}
	
	// Return success response as specified
	response := map[string]string{
		"message": "Request successful",
	}
	
	// Encode and send the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, log the error
		log.Printf("Error encoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

