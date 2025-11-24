package main

import (
	"sync"
	"time"
)

// RateLimiter manages rate limiting for IP addresses
// It tracks the number of requests per IP and enforces a limit
type RateLimiter struct {
	// Requests maps IP addresses to their request count in the current window
	Requests map[string]int
	
	// Mutex ensures thread-safe access to the Requests map
	// This prevents race conditions when multiple goroutines access the map concurrently
	Mutex sync.Mutex
	
	// MaxRequests is the maximum number of requests allowed per time window
	MaxRequests int
	
	// WindowDuration is the time window for rate limiting (e.g., 1 minute)
	WindowDuration time.Duration
}

// NewRateLimiter creates a new RateLimiter instance
// maxRequests: maximum requests allowed per window (e.g., 5)
// windowDuration: time window duration (e.g., 1 minute)
func NewRateLimiter(maxRequests int, windowDuration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		Requests:       make(map[string]int),
		MaxRequests:    maxRequests,
		WindowDuration: windowDuration,
	}
	
	// Start the background goroutine that resets request counts every minute
	rl.startResetTicker()
	
	return rl
}

// startResetTicker starts a goroutine that periodically resets the request counts
// This runs in the background and clears the Requests map every time window
// This is similar to how services like Stripe/OpenAI reset their rate limits
func (rl *RateLimiter) startResetTicker() {
	ticker := time.NewTicker(rl.WindowDuration)
	
	go func() {
		for range ticker.C {
			// Lock the mutex before modifying the map
			rl.Mutex.Lock()
			
			// Clear all request counts - this effectively resets the rate limit window
			// All IPs get a fresh start for the next minute
			rl.Requests = make(map[string]int)
			
			// Unlock after modification
			rl.Mutex.Unlock()
		}
	}()
}

// AllowRequest checks if a request from the given IP should be allowed
// Returns true if the request is allowed, false if rate limit is exceeded
func (rl *RateLimiter) AllowRequest(ip string) bool {
	// Lock the mutex to ensure thread-safe access
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()
	
	// Get current request count for this IP (defaults to 0 if IP not in map)
	currentCount := rl.Requests[ip]
	
	// Check if the IP has exceeded the maximum requests
	if currentCount >= rl.MaxRequests {
		// Rate limit exceeded - don't increment, just return false
		return false
	}
	
	// Increment the request count for this IP
	rl.Requests[ip] = currentCount + 1
	
	// Request is allowed
	return true
}

// GetRemainingRequests returns the number of remaining requests for an IP
// This is useful for debugging or returning rate limit headers
func (rl *RateLimiter) GetRemainingRequests(ip string) int {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()
	
	currentCount := rl.Requests[ip]
	remaining := rl.MaxRequests - currentCount
	
	if remaining < 0 {
		return 0
	}
	
	return remaining
}

