package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// RateLimitMiddleware creates a middleware function that applies rate limiting
// This middleware runs before the actual handler and checks if the request should be allowed
func RateLimitMiddleware(limiter *RateLimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Extract the client's IP address from the request
			// This handles various proxy scenarios (X-Forwarded-For, X-Real-IP)
			ip := getClientIP(r)
			
			// Check if the request should be allowed based on rate limiting
			if !limiter.AllowRequest(ip) {
				// Rate limit exceeded - return 429 Too Many Requests
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				
				// Return JSON error response as specified
				errorResponse := map[string]string{
					"error": "Rate limit exceeded",
				}
				
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
			
			// Request is allowed - proceed to the next handler
			next(w, r)
		}
	}
}

// getClientIP extracts the client's IP address from the HTTP request
// It checks various headers that proxies/load balancers might set
// Falls back to RemoteAddr if no headers are present
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by proxies/load balancers)
	// This header can contain multiple IPs, so we take the first one
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs separated by commas
		// The first IP is usually the original client IP
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check X-Real-IP header (alternative header used by some proxies)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return strings.TrimSpace(realIP)
	}
	
	// Fall back to RemoteAddr (direct connection)
	// RemoteAddr includes port, so we extract just the IP
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	
	return ip
}

