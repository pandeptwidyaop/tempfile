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
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(store Store, ipDetector IPDetector, config *Config) RateLimiter {
	return &rateLimiter{
		store:            store,
		ipDetector:       ipDetector,
		uploadsPerMinute: config.UploadsPerMinute,
		bytesPerHour:     config.BytesPerHour,
		windowMinutes:    config.WindowMinutes,
	}
}

// CheckLimits verifies if an IP can upload a file of the given size
func (r *rateLimiter) CheckLimits(ip string, fileSize int64) (*LimitStatus, error) {
	now := time.Now()
	
	// Calculate time windows
	uploadWindow := time.Duration(r.windowMinutes) * time.Minute
	bytesWindow := time.Hour // Always use 1 hour for bytes limit
	
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
		IsLimited:    false,
	}
	
	// Check upload count limit
	if uploadCount >= r.uploadsPerMinute {
		status.IsLimited = true
		status.LimitReason = fmt.Sprintf("Upload limit: %d uploads per %d minutes exceeded", 
			r.uploadsPerMinute, r.windowMinutes)
		
		return status, NewRateLimitError(
			ip,
			"upload_count",
			status.LimitReason,
			r.calculateRetryAfter(uploadWindow),
			map[string]interface{}{
				"uploads_used":  uploadCount,
				"uploads_limit": r.uploadsPerMinute,
				"window_minutes": r.windowMinutes,
			},
		)
	}
	
	// Check bytes limit (including the new file)
	if bytesUsed+fileSize > r.bytesPerHour {
		status.IsLimited = true
		status.LimitReason = fmt.Sprintf("Bytes limit: %d bytes per hour exceeded", r.bytesPerHour)
		
		return status, NewRateLimitError(
			ip,
			"bytes_limit",
			status.LimitReason,
			r.calculateRetryAfter(bytesWindow),
			map[string]interface{}{
				"bytes_used":     bytesUsed,
				"bytes_limit":    r.bytesPerHour,
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
	
	// Create IP detector
	ipDetector := NewIPDetector(config.TrustedProxies, config.IPHeaders)
	
	return NewRateLimiter(store, ipDetector, config)
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