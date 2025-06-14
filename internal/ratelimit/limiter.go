package ratelimit

import (
	"fmt"
	"time"
)

// rateLimiter implements the RateLimiter interface
type rateLimiter struct {
	store            Store
	ipDetector       IPDetector
	uploadsPerMinute int
	bytesPerHour     int64
	windowMinutes    int
	customLimits     map[string]EndpointConfig
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(store Store, ipDetector IPDetector, config *Config) RateLimiter {
	return &rateLimiter{
		store:            store,
		ipDetector:       ipDetector,
		uploadsPerMinute: config.UploadsPerMinute,
		bytesPerHour:     config.BytesPerHour,
		windowMinutes:    config.WindowMinutes,
		customLimits:     config.CustomLimits,
	}
}

// CheckLimits verifies if an IP can upload a file of the given size
func (r *rateLimiter) CheckLimits(ip string, fileSize int64) (*LimitStatus, error) {
	return r.CheckLimitsForEndpoint(ip, fileSize, "")
}

// CheckLimitsForEndpoint verifies if an IP can upload a file with endpoint-specific limits
func (r *rateLimiter) CheckLimitsForEndpoint(ip string, fileSize int64, endpoint string) (*LimitStatus, error) {
	// Check if IP is whitelisted
	if detector, ok := r.ipDetector.(*ipDetector); ok {
		if detector.IsWhitelisted(ip) {
			// Return unlimited status for whitelisted IPs
			return &LimitStatus{
				IP:           ip,
				UploadsUsed:  0,
				UploadsLimit: -1, // -1 indicates unlimited
				BytesUsed:    0,
				BytesLimit:   -1,
				WindowStart:  time.Now(),
				WindowEnd:    time.Now(),
				ResetTime:    time.Now(),
				IsLimited:    false,
				LimitReason:  "whitelisted",
			}, nil
		}
	}

	// Get limits for this endpoint
	uploadsPerMinute := r.uploadsPerMinute
	bytesPerHour := r.bytesPerHour
	windowMinutes := r.windowMinutes

	if endpoint != "" && r.customLimits != nil {
		if customLimit, exists := r.customLimits[endpoint]; exists {
			uploadsPerMinute = customLimit.UploadsPerMinute
			bytesPerHour = customLimit.BytesPerHour
			windowMinutes = customLimit.WindowMinutes
		}
	}

	now := time.Now()
	uploadWindow := time.Duration(windowMinutes) * time.Minute
	bytesWindow := time.Hour // Always use 1 hour for bytes limit

	// Try atomic operation for Redis
	if atomicStore, ok := r.store.(AtomicStore); ok {
		allowed, uploadCount, bytesUsed, reason, err := atomicStore.AtomicCheckAndIncrement(
			ip, fileSize, uploadWindow, uploadsPerMinute, bytesPerHour)

		if err != nil {
			return nil, fmt.Errorf("atomic rate limit check failed: %w", err)
		}

		status := &LimitStatus{
			IP:           ip,
			UploadsUsed:  uploadCount,
			UploadsLimit: uploadsPerMinute,
			BytesUsed:    bytesUsed,
			BytesLimit:   bytesPerHour,
			WindowStart:  now.Add(-uploadWindow),
			WindowEnd:    now,
			ResetTime:    now.Add(uploadWindow),
			IsLimited:    !allowed,
		}

		if !allowed {
			if reason == "upload_limit" {
				status.LimitReason = fmt.Sprintf("Upload limit: %d uploads per %d minutes exceeded",
					uploadsPerMinute, windowMinutes)
			} else {
				status.LimitReason = fmt.Sprintf("Bytes limit: %d bytes per hour exceeded", bytesPerHour)
			}

			return status, NewRateLimitError(
				ip,
				reason,
				status.LimitReason,
				r.calculateRetryAfter(uploadWindow),
				map[string]interface{}{
					"uploads_used":   uploadCount,
					"uploads_limit":  uploadsPerMinute,
					"bytes_used":     bytesUsed,
					"bytes_limit":    bytesPerHour,
					"window_minutes": windowMinutes,
				},
			)
		}

		return status, nil
	}

	// Fallback to non-atomic operations for memory store
	uploadCount, err := r.store.GetUploadCount(ip, uploadWindow)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload count: %w", err)
	}

	bytesUsed, err := r.store.GetBytesUsed(ip, bytesWindow)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes used: %w", err)
	}

	// Create status
	status := &LimitStatus{
		IP:           ip,
		UploadsUsed:  uploadCount,
		UploadsLimit: uploadsPerMinute,
		BytesUsed:    bytesUsed,
		BytesLimit:   bytesPerHour,
		WindowStart:  now.Add(-uploadWindow),
		WindowEnd:    now,
		ResetTime:    now.Add(uploadWindow),
		IsLimited:    false,
	}

	// Check upload count limit
	if uploadCount >= uploadsPerMinute {
		status.IsLimited = true
		status.LimitReason = fmt.Sprintf("Upload limit: %d uploads per %d minutes exceeded",
			uploadsPerMinute, windowMinutes)

		return status, NewRateLimitError(
			ip,
			"upload_count",
			status.LimitReason,
			r.calculateRetryAfter(uploadWindow),
			map[string]interface{}{
				"uploads_used":   uploadCount,
				"uploads_limit":  uploadsPerMinute,
				"window_minutes": windowMinutes,
			},
		)
	}

	// Check bytes limit (including the new file)
	if bytesUsed+fileSize > bytesPerHour {
		status.IsLimited = true
		status.LimitReason = fmt.Sprintf("Bytes limit: %d bytes per hour exceeded", bytesPerHour)

		return status, NewRateLimitError(
			ip,
			"bytes_limit",
			status.LimitReason,
			r.calculateRetryAfter(bytesWindow),
			map[string]interface{}{
				"bytes_used":     bytesUsed,
				"bytes_limit":    bytesPerHour,
				"file_size":      fileSize,
				"total_would_be": bytesUsed + fileSize,
			},
		)
	}

	return status, nil
}

// UpdateCounters increments the counters after a successful upload
func (r *rateLimiter) UpdateCounters(ip string, fileSize int64) error {
	uploadWindow := time.Duration(r.windowMinutes) * time.Minute

	return r.store.IncrementUpload(ip, fileSize, uploadWindow)
}

// GetStatus returns the current rate limit status for an IP
func (r *rateLimiter) GetStatus(ip string) (*LimitStatus, error) {
	now := time.Now()

	// Calculate time windows
	uploadWindow := time.Duration(r.windowMinutes) * time.Minute
	bytesWindow := time.Hour

	// Get current usage
	uploadCount, err := r.store.GetUploadCount(ip, uploadWindow)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload count: %w", err)
	}

	bytesUsed, err := r.store.GetBytesUsed(ip, bytesWindow)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes used: %w", err)
	}

	// Create status
	status := &LimitStatus{
		IP:           ip,
		UploadsUsed:  uploadCount,
		UploadsLimit: r.uploadsPerMinute,
		BytesUsed:    bytesUsed,
		BytesLimit:   r.bytesPerHour,
		WindowStart:  now.Add(-uploadWindow),
		WindowEnd:    now,
		ResetTime:    now.Add(uploadWindow),
		IsLimited:    uploadCount >= r.uploadsPerMinute || bytesUsed >= r.bytesPerHour,
	}

	if status.IsLimited {
		if uploadCount >= r.uploadsPerMinute {
			status.LimitReason = "Upload count limit exceeded"
		} else {
			status.LimitReason = "Bytes limit exceeded"
		}
	}

	return status, nil
}

// Close closes the rate limiter and releases resources
func (r *rateLimiter) Close() error {
	if r.store != nil {
		return r.store.Close()
	}
	return nil
}

// calculateRetryAfter calculates when the client should retry (in seconds)
func (r *rateLimiter) calculateRetryAfter(window time.Duration) int {
	// Return the window duration in seconds as a conservative estimate
	return int(window.Seconds())
}

// NewDefaultMemoryRateLimiter creates a rate limiter with in-memory storage and default settings
func NewDefaultMemoryRateLimiter(config *Config) RateLimiter {
	// Create memory store with default settings
	store := NewMemoryStore(10000, 5*time.Minute) // 10k IPs, cleanup every 5 minutes

	// Create IP detector with whitelist support
	ipDetector := NewIPDetectorWithWhitelist(config.TrustedProxies, config.IPHeaders, config.WhitelistIPs)

	return NewRateLimiter(store, ipDetector, config)
}

// NewRedisRateLimiter creates a rate limiter with Redis storage
func NewRedisRateLimiter(config *Config) (RateLimiter, error) {
	// Create Redis store
	store, err := NewRedisStore(config.RedisURL, config.RedisPassword, config.RedisDB, config.RedisPoolSize, config.RedisTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis store: %w", err)
	}

	// Create IP detector with whitelist support
	ipDetector := NewIPDetectorWithWhitelist(config.TrustedProxies, config.IPHeaders, config.WhitelistIPs)

	return NewRateLimiter(store, ipDetector, config), nil
}

// ValidateConfig validates the rate limiter configuration
func ValidateConfig(config *Config) error {
	if config.UploadsPerMinute <= 0 {
		return fmt.Errorf("uploads per minute must be positive, got %d", config.UploadsPerMinute)
	}

	if config.BytesPerHour <= 0 {
		return fmt.Errorf("bytes per hour must be positive, got %d", config.BytesPerHour)
	}

	if config.WindowMinutes <= 0 {
		return fmt.Errorf("window minutes must be positive, got %d", config.WindowMinutes)
	}

	if config.Store != "memory" && config.Store != "redis" {
		return fmt.Errorf("store must be 'memory' or 'redis', got '%s'", config.Store)
	}

	return nil
}
