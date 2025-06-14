package ratelimit

import "errors"

// Common rate limiting errors
var (
	// ErrStoreClosed indicates the store has been closed
	ErrStoreClosed = errors.New("rate limit store is closed")

	// ErrStoreCapacityExceeded indicates the store has reached maximum capacity
	ErrStoreCapacityExceeded = errors.New("rate limit store capacity exceeded")

	// ErrRateLimitExceeded indicates a rate limit has been exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrInvalidIP indicates an invalid IP address was provided
	ErrInvalidIP = errors.New("invalid IP address")

	// ErrInvalidConfiguration indicates invalid configuration
	ErrInvalidConfiguration = errors.New("invalid rate limiter configuration")

	// ErrRedisConnection indicates a Redis connection error
	ErrRedisConnection = errors.New("Redis connection error")

	// ErrRedisOperation indicates a Redis operation error
	ErrRedisOperation = errors.New("Redis operation error")
)

// RateLimitError represents a rate limit exceeded error with details
type RateLimitError struct {
	IP           string
	LimitType    string
	Message      string
	RetryAfter   int
	CurrentUsage map[string]interface{}
}

func (e *RateLimitError) Error() string {
	return e.Message
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(ip, limitType, message string, retryAfter int, usage map[string]interface{}) *RateLimitError {
	return &RateLimitError{
		IP:           ip,
		LimitType:    limitType,
		Message:      message,
		RetryAfter:   retryAfter,
		CurrentUsage: usage,
	}
}
