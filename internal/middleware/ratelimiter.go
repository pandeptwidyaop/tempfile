package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/tempfile/internal/ratelimit"
)

// RateLimiterConfig holds the configuration for the rate limiter middleware
type RateLimiterConfig struct {
	// RateLimiter is the rate limiter instance to use
	RateLimiter ratelimit.RateLimiter
	
	// IPDetector is the IP detector instance to use
	IPDetector ratelimit.IPDetector
	
	// SkipPaths are paths that should skip rate limiting
	SkipPaths []string
	
	// SkipSuccessfulRequests if true, only failed requests count towards rate limit
	SkipSuccessfulRequests bool
	
	// KeyGenerator allows custom key generation for rate limiting
	KeyGenerator func(c *fiber.Ctx) string
}

// New creates a new rate limiter middleware
func NewRateLimiter(config RateLimiterConfig) fiber.Handler {
	// Set defaults
	if config.KeyGenerator == nil {
		config.KeyGenerator = defaultKeyGenerator(config.IPDetector)
	}
	
	if config.SkipPaths == nil {
		config.SkipPaths = []string{"/health"}
	}
	
	return func(c *fiber.Ctx) error {
		// Skip rate limiting for certain paths
		path := c.Path()
		for _, skipPath := range config.SkipPaths {
			if path == skipPath {
				return c.Next()
			}
		}
		
		// Generate key (usually IP address)
		key := config.KeyGenerator(c)
		if key == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Unable to determine client identifier",
				"code":  "IP_DETECTION_FAILED",
			})
		}
		
		// Get estimated file size for pre-validation
		fileSize := getEstimatedFileSize(c)
		
		// Check rate limits
		status, err := config.RateLimiter.CheckLimits(key, fileSize)
		if err != nil {
			// Check if it's a rate limit error
			if rateLimitErr, ok := err.(*ratelimit.RateLimitError); ok {
				return handleRateLimitExceeded(c, rateLimitErr, status)
			}
			
			// Other errors (store errors, etc.)
			return c.Status(500).JSON(fiber.Map{
				"error": "Rate limit check failed",
				"code":  "RATE_LIMIT_ERROR",
			})
		}
		
		// Store information for post-processing
		c.Locals("rate_limit_key", key)
		c.Locals("rate_limit_checked", true)
		c.Locals("rate_limit_estimated_size", fileSize)
		
		// Add rate limit headers to response
		addRateLimitHeaders(c, status)
		
		return c.Next()
	}
}

// PostProcess creates a middleware that updates counters after successful upload
func NewRateLimiterPostProcess(rateLimiter ratelimit.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only update if rate limit was checked and request was successful
		if c.Locals("rate_limit_checked") == true && c.Response().StatusCode() < 400 {
			key := c.Locals("rate_limit_key")
			if key != nil {
				// Get actual file size if available, otherwise use estimated
				actualSize := getActualFileSize(c)
				if actualSize == 0 {
					if estimatedSize := c.Locals("rate_limit_estimated_size"); estimatedSize != nil {
						actualSize = estimatedSize.(int64)
					}
				}
				
				// Update counters
				if err := rateLimiter.UpdateCounters(key.(string), actualSize); err != nil {
					// Log error but don't fail the request
					// TODO: Add proper logging
				}
			}
		}
		
		return c.Next()
	}
}

// defaultKeyGenerator creates a default key generator using IP detection
func defaultKeyGenerator(ipDetector ratelimit.IPDetector) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		// Extract headers
		headers := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})
		
		// Get remote address
		remoteAddr := c.Context().RemoteAddr().String()
		
		// Detect real IP
		return ipDetector.GetRealIP(headers, remoteAddr)
	}
}

// getEstimatedFileSize attempts to get file size from Content-Length header
func getEstimatedFileSize(c *fiber.Ctx) int64 {
	// Try Content-Length header first
	if contentLength := c.Get("Content-Length"); contentLength != "" {
		if size, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
			return size
		}
	}
	
	// Try to get from multipart form (this might parse the form)
	if form, err := c.MultipartForm(); err == nil && form != nil {
		if files := form.File["file"]; len(files) > 0 {
			return files[0].Size
		}
	}
	
	// Default to 0 if we can't determine size
	return 0
}

// getActualFileSize gets the actual file size after processing
func getActualFileSize(c *fiber.Ctx) int64 {
	// Check if actual size was stored during processing
	if actualSize := c.Locals("actual_file_size"); actualSize != nil {
		if size, ok := actualSize.(int64); ok {
			return size
		}
	}
	
	// Try to get from multipart form again
	if form, err := c.MultipartForm(); err == nil && form != nil {
		if files := form.File["file"]; len(files) > 0 {
			return files[0].Size
		}
	}
	
	return 0
}

// handleRateLimitExceeded handles rate limit exceeded errors
func handleRateLimitExceeded(c *fiber.Ctx, rateLimitErr *ratelimit.RateLimitError, status *ratelimit.LimitStatus) error {
	// Set Retry-After header
	c.Set("Retry-After", strconv.Itoa(rateLimitErr.RetryAfter))
	
	// Determine if this is a web request or API request
	acceptHeader := c.Get("Accept")
	isWebRequest := strings.Contains(acceptHeader, "text/html")
	
	if isWebRequest {
		// For web requests, return HTML error page
		return c.Status(429).SendString(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Rate Limit Exceeded</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .error { color: #d32f2f; }
        .info { color: #1976d2; margin-top: 20px; }
    </style>
</head>
<body>
    <h1 class="error">Rate Limit Exceeded</h1>
    <p>%s</p>
    <div class="info">
        <p><strong>Current Usage:</strong></p>
        <ul>
            <li>Uploads: %d / %d</li>
            <li>Data: %.2f MB / %.2f MB</li>
        </ul>
        <p>Please try again in %d seconds.</p>
    </div>
</body>
</html>`,
			rateLimitErr.Message,
			status.UploadsUsed, status.UploadsLimit,
			float64(status.BytesUsed)/(1024*1024), float64(status.BytesLimit)/(1024*1024),
			rateLimitErr.RetryAfter,
		))
	}
	
	// For API requests, return JSON
	return c.Status(429).JSON(fiber.Map{
		"error":   "Rate limit exceeded",
		"code":    "RATE_LIMIT_EXCEEDED",
		"message": rateLimitErr.Message,
		"details": fiber.Map{
			"limit_type": rateLimitErr.LimitType,
			"reason":     rateLimitErr.Message,
		},
		"current_usage": fiber.Map{
			"ip":            status.IP,
			"uploads_used":  status.UploadsUsed,
			"uploads_limit": status.UploadsLimit,
			"bytes_used":    status.BytesUsed,
			"bytes_limit":   status.BytesLimit,
			"window_start":  status.WindowStart.Format("2006-01-02T15:04:05Z07:00"),
			"window_end":    status.WindowEnd.Format("2006-01-02T15:04:05Z07:00"),
		},
		"retry_after": rateLimitErr.RetryAfter,
		"reset_time":  status.ResetTime.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// addRateLimitHeaders adds rate limit information to response headers
func addRateLimitHeaders(c *fiber.Ctx, status *ratelimit.LimitStatus) {
	c.Set("X-RateLimit-Limit-Uploads", strconv.Itoa(status.UploadsLimit))
	c.Set("X-RateLimit-Remaining-Uploads", strconv.Itoa(status.UploadsLimit-status.UploadsUsed))
	c.Set("X-RateLimit-Limit-Bytes", strconv.FormatInt(status.BytesLimit, 10))
	c.Set("X-RateLimit-Remaining-Bytes", strconv.FormatInt(status.BytesLimit-status.BytesUsed, 10))
	c.Set("X-RateLimit-Reset", strconv.FormatInt(status.ResetTime.Unix(), 10))
}