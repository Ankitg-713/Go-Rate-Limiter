# Go Rate Limiter API

A production-ready rate limiting API server built with Go, similar to OpenAI/Stripe rate limiting. This implementation uses only standard Go libraries and provides IP-based rate limiting with automatic reset windows.

## Features

- **IP-based rate limiting**: Tracks requests per IP address
- **Configurable limits**: 5 requests per minute (configurable)
- **Automatic reset**: Rate limits reset every minute using background goroutines
- **Thread-safe**: Uses `sync.Mutex` to prevent race conditions
- **Standard library only**: Built with `net/http` and standard Go packages
- **RESTful API**: Simple JSON API with clear error responses

## Project Structure

```
go-rate-limiter/
├── main.go           # HTTP server and endpoint handlers
├── rate_limiter.go   # Rate limiter logic with goroutine-based reset
├── middleware.go      # Middleware for applying rate limits
├── go.mod            # Go module file
└── README.md         # This file
```

## How It Works

1. **Rate Limiter**: Maintains a map of IP addresses to request counts
2. **Background Reset**: A goroutine with `time.Ticker` resets all counts every minute
3. **Middleware**: Intercepts requests, checks rate limits, and blocks if exceeded
4. **Thread Safety**: `sync.Mutex` ensures concurrent requests don't cause race conditions

## Setup and Installation

### Prerequisites

- Go 1.21 or higher installed on your system
- Terminal/command line access

### Installation Steps

1. **Navigate to the project directory**:
   ```bash
   cd go-rate-limiter
   ```

2. **Initialize the Go module** (if not already done):
   ```bash
   go mod init go-rate-limiter
   ```

3. **Install dependencies** (this project uses only standard library, so no external packages needed):
   ```bash
   go mod tidy
   ```

4. **Run the server**:
   ```bash
   go run main.go rate_limiter.go middleware.go
   ```

   Or simply:
   ```bash
   go run .
   ```

The server will start on `http://localhost:8080`

## Usage

### Making Requests

**Successful request** (within rate limit):
```bash
curl http://localhost:8080/api/data
```

Response:
```json
{"message": "Request successful"}
```

**Rate limit exceeded** (after 5 requests):
```bash
curl http://localhost:8080/api/data
```

Response (HTTP 429):
```json
{"error": "Rate limit exceeded"}
```

### Testing Rate Limits

You can test the rate limiter by making multiple rapid requests:

```bash
# Make 6 requests quickly (5th should succeed, 6th should fail)
for i in {1..6}; do
  echo "Request $i:"
  curl http://localhost:8080/api/data
  echo ""
done
```

After 1 minute, the rate limit resets automatically and you can make 5 more requests.

## Configuration

You can modify the rate limiting parameters in `main.go`:

```go
const (
    maxRequestsPerMinute = 5        // Change this to adjust the limit
    rateLimitWindow      = 1 * time.Minute  // Change this to adjust the time window
)
```

## API Endpoints

### GET /api/data

Returns a success message if the rate limit is not exceeded.

**Response (200 OK)**:
```json
{
  "message": "Request successful"
}
```

**Response (429 Too Many Requests)**:
```json
{
  "error": "Rate limit exceeded"
}
```

## Architecture Details

### Rate Limiter (`rate_limiter.go`)

- **Struct**: Contains a map of IPs to request counts, a mutex for thread safety, and configuration
- **Reset Mechanism**: Background goroutine uses `time.Ticker` to reset counts every minute
- **Thread Safety**: All map operations are protected by `sync.Mutex`

### Middleware (`middleware.go`)

- **IP Extraction**: Handles various proxy headers (X-Forwarded-For, X-Real-IP)
- **Rate Check**: Calls the rate limiter before allowing requests through
- **Error Response**: Returns 429 status with JSON error message when limit exceeded

### Main Server (`main.go`)

- **HTTP Server**: Uses standard `net/http` package
- **Route Registration**: Applies middleware to the `/api/data` endpoint
- **Handler**: Simple JSON response handler

## Interview Explanation Points

1. **Why maps?**: Fast O(1) lookup for IP addresses
2. **Why goroutines?**: Non-blocking background reset that doesn't affect request handling
3. **Why time.Ticker?**: Precise periodic reset without manual intervention
4. **Why sync.Mutex?**: Prevents race conditions when multiple goroutines access the map concurrently
5. **Why middleware?**: Clean separation of concerns, reusable rate limiting logic

## Production Considerations

For production use, consider:

- **Distributed rate limiting**: Use Redis or similar for multi-server deployments
- **Persistence**: Save rate limit state across server restarts
- **Rate limit headers**: Add `X-RateLimit-Remaining` headers to responses
- **Logging**: Add structured logging for monitoring
- **Metrics**: Track rate limit hits and misses
- **Configuration**: Load limits from environment variables or config files

## License

This project is provided as-is for educational and interview purposes.

