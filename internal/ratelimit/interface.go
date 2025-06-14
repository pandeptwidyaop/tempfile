package ratelimit

import (
	"time"
)

// LimitStatus represents the current rate limit status for an IP
type LimitStatus struct {
	IP           string    `json:"ip"`
	UploadsUsed  int       `json:"uploads_used"`
	UploadsLimit int       `json:"uploads_limit"`
	BytesUsed    int64     `json:"bytes_used"`
	BytesLimit   int64     `json:"bytes_limit"`
	WindowStart  time.Time `json:"window_start"`
	WindowEnd    time.Time `json:"window_end"`
	ResetTime    time.Time `json:"reset_time"`
	IsLimited    bool      `json:"is_limited"`
	LimitReason  string    `json:"limit_reason,omitempty"`
}

// Store interface defines the storage backend for rate limiting data
type Store interface {
	// GetUploadCount returns the number of uploads for an IP within the time window
	GetUploadCount(ip string, window time.Duration) (int, error)

	// GetBytesUsed returns the total bytes uploaded for an IP within the time window
	GetBytesUsed(ip string, window time.Duration) (int64, error)

	// IncrementUpload records a new upload for an IP with the given file size
	IncrementUpload(ip string, fileSize int64, window time.Duration) error

	// Cleanup removes expired entries from the store
	Cleanup() error

	// HealthCheck verifies the store is functioning properly
	HealthCheck() error

	// Close closes the store and releases resources
	Close() error
}

// RateLimiter interface defines the main rate limiting functionality
type RateLimiter interface {
	// CheckLimits verifies if an IP can upload a file of the given size
	CheckLimits(ip string, fileSize int64) (*LimitStatus, error)

	// CheckLimitsForEndpoint verifies if an IP can upload a file with endpoint-specific limits
	CheckLimitsForEndpoint(ip string, fileSize int64, endpoint string) (*LimitStatus, error)

	// UpdateCounters increments the counters after a successful upload
	UpdateCounters(ip string, fileSize int64) error

	// GetStatus returns the current rate limit status for an IP
	GetStatus(ip string) (*LimitStatus, error)

	// Close closes the rate limiter and releases resources
	Close() error
}

// IPDetector interface defines IP detection functionality
type IPDetector interface {
	// GetRealIP extracts the real client IP from HTTP headers
	GetRealIP(headers map[string]string, remoteAddr string) string

	// IsTrustedProxy checks if an IP is a trusted proxy
	IsTrustedProxy(ip string) bool
}

// Config holds rate limiter configuration
type Config struct {
	Store            string
	UploadsPerMinute int
	BytesPerHour     int64
	WindowMinutes    int
	TrustedProxies   []string
	IPHeaders        []string
	WhitelistIPs     []string
	CustomLimits     map[string]EndpointConfig
	RedisURL         string
	RedisPassword    string
	RedisDB          int
	RedisPoolSize    int
	RedisTimeout     int
}

// EndpointConfig holds custom rate limits for specific endpoints
type EndpointConfig struct {
	UploadsPerMinute int
	BytesPerHour     int64
	WindowMinutes    int
}

// AtomicStore extends Store with atomic operations for Redis
type AtomicStore interface {
	Store
	// AtomicCheckAndIncrement performs atomic rate limit check and increment
	AtomicCheckAndIncrement(ip string, fileSize int64, window time.Duration, uploadLimit int, bytesLimit int64) (bool, int, int64, string, error)
}
