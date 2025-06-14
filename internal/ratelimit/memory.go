package ratelimit

import (
	"sync"
	"time"
)

// UploadRecord represents a single upload record
type UploadRecord struct {
	Timestamp time.Time
	FileSize  int64
}

// memoryStore implements the Store interface using in-memory storage
type memoryStore struct {
	mu          sync.RWMutex
	uploads     map[string][]UploadRecord
	maxEntries  int
	cleanupTick time.Duration
	stopCleanup chan struct{}
	closed      bool
}

// NewMemoryStore creates a new in-memory rate limit store
func NewMemoryStore(maxEntries int, cleanupInterval time.Duration) Store {
	store := &memoryStore{
		uploads:     make(map[string][]UploadRecord),
		maxEntries:  maxEntries,
		cleanupTick: cleanupInterval,
		stopCleanup: make(chan struct{}),
	}
	
	// Start cleanup goroutine
	go store.cleanupLoop()
	
	return store
}

// GetUploadCount returns the number of uploads for an IP within the time window
func (s *memoryStore) GetUploadCount(ip string, window time.Duration) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if s.closed {
		return 0, ErrStoreClosed
	}
	
	records, exists := s.uploads[ip]
	if !exists {
		return 0, nil
	}
	
	cutoff := time.Now().Add(-window)
	count := 0
	
	for _, record := range records {
		if record.Timestamp.After(cutoff) {
			count++
		}
	}
	
	return count, nil
}

// GetBytesUsed returns the total bytes uploaded for an IP within the time window
func (s *memoryStore) GetBytesUsed(ip string, window time.Duration) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if s.closed {
		return 0, ErrStoreClosed
	}
	
	records, exists := s.uploads[ip]
	if !exists {
		return 0, nil
	}
	
	cutoff := time.Now().Add(-window)
	var totalBytes int64
	
	for _, record := range records {
		if record.Timestamp.After(cutoff) {
			totalBytes += record.FileSize
		}
	}
	
	return totalBytes, nil
}

// IncrementUpload records a new upload for an IP with the given file size
func (s *memoryStore) IncrementUpload(ip string, fileSize int64, window time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.closed {
		return ErrStoreClosed
	}
	
	// Check if we're at max capacity
	if len(s.uploads) >= s.maxEntries {
		// Try to clean up first
		s.cleanupExpiredEntries(window)
		
		// If still at capacity, reject
		if len(s.uploads) >= s.maxEntries {
			return ErrStoreCapacityExceeded
		}
	}
	
	now := time.Now()
	record := UploadRecord{
		Timestamp: now,
		FileSize:  fileSize,
	}
	
	s.uploads[ip] = append(s.uploads[ip], record)
	
	return nil
}

// Cleanup removes expired entries from the store
func (s *memoryStore) Cleanup() error {
	return s.CleanupWithWindow(24 * time.Hour)
}

// CleanupWithWindow removes expired entries using a specific window
func (s *memoryStore) CleanupWithWindow(window time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.closed {
		return ErrStoreClosed
	}
	
	s.cleanupExpiredEntries(window)
	
	return nil
}

// cleanupExpiredEntries removes expired entries (must be called with lock held)
func (s *memoryStore) cleanupExpiredEntries(window time.Duration) {
	cutoff := time.Now().Add(-window)
	
	for ip, records := range s.uploads {
		var validRecords []UploadRecord
		
		for _, record := range records {
			if record.Timestamp.After(cutoff) {
				validRecords = append(validRecords, record)
			}
		}
		
		if len(validRecords) == 0 {
			delete(s.uploads, ip)
		} else {
			s.uploads[ip] = validRecords
		}
	}
}

// cleanupLoop runs periodic cleanup
func (s *memoryStore) cleanupLoop() {
	ticker := time.NewTicker(s.cleanupTick)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.Cleanup()
		case <-s.stopCleanup:
			return
		}
	}
}

// HealthCheck verifies the store is functioning properly
func (s *memoryStore) HealthCheck() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if s.closed {
		return ErrStoreClosed
	}
	
	return nil
}

// Close closes the store and releases resources
func (s *memoryStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.closed {
		return nil
	}
	
	s.closed = true
	close(s.stopCleanup)
	s.uploads = nil
	
	return nil
}

// GetStats returns statistics about the memory store
func (s *memoryStore) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	totalRecords := 0
	for _, records := range s.uploads {
		totalRecords += len(records)
	}
	
	return map[string]interface{}{
		"type":          "memory",
		"active_ips":    len(s.uploads),
		"total_records": totalRecords,
		"max_entries":   s.maxEntries,
		"closed":        s.closed,
	}
}