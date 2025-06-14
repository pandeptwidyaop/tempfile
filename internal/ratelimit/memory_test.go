package ratelimit

import (
	"testing"
	"time"
)

func TestMemoryStore_GetUploadCount(t *testing.T) {
	store := NewMemoryStore(100, time.Minute)
	defer store.Close()

	ip := "203.0.113.1"
	window := time.Minute

	// Initially should be 0
	count, err := store.GetUploadCount(ip, window)
	if err != nil {
		t.Fatalf("GetUploadCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetUploadCount() = %v, want 0", count)
	}

	// Add some uploads
	err = store.IncrementUpload(ip, 1024, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	err = store.IncrementUpload(ip, 2048, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	// Should now be 2
	count, err = store.GetUploadCount(ip, window)
	if err != nil {
		t.Fatalf("GetUploadCount() error = %v", err)
	}
	if count != 2 {
		t.Errorf("GetUploadCount() = %v, want 2", count)
	}
}

func TestMemoryStore_GetBytesUsed(t *testing.T) {
	store := NewMemoryStore(100, time.Minute)
	defer store.Close()

	ip := "203.0.113.1"
	window := time.Minute

	// Initially should be 0
	bytes, err := store.GetBytesUsed(ip, window)
	if err != nil {
		t.Fatalf("GetBytesUsed() error = %v", err)
	}
	if bytes != 0 {
		t.Errorf("GetBytesUsed() = %v, want 0", bytes)
	}

	// Add some uploads
	err = store.IncrementUpload(ip, 1024, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	err = store.IncrementUpload(ip, 2048, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	// Should now be 3072 (1024 + 2048)
	bytes, err = store.GetBytesUsed(ip, window)
	if err != nil {
		t.Fatalf("GetBytesUsed() error = %v", err)
	}
	if bytes != 3072 {
		t.Errorf("GetBytesUsed() = %v, want 3072", bytes)
	}
}

func TestMemoryStore_SlidingWindow(t *testing.T) {
	store := NewMemoryStore(100, time.Minute)
	defer store.Close()

	ip := "203.0.113.1"
	window := 100 * time.Millisecond

	// Add upload
	err := store.IncrementUpload(ip, 1024, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	// Should be visible immediately
	count, err := store.GetUploadCount(ip, window)
	if err != nil {
		t.Fatalf("GetUploadCount() error = %v", err)
	}
	if count != 1 {
		t.Errorf("GetUploadCount() = %v, want 1", count)
	}

	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)

	// Should now be 0 (outside window)
	count, err = store.GetUploadCount(ip, window)
	if err != nil {
		t.Fatalf("GetUploadCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetUploadCount() = %v, want 0", count)
	}
}

func TestMemoryStore_MultipleIPs(t *testing.T) {
	store := NewMemoryStore(100, time.Minute)
	defer store.Close()

	ip1 := "203.0.113.1"
	ip2 := "203.0.113.2"
	window := time.Minute

	// Add uploads for different IPs
	err := store.IncrementUpload(ip1, 1024, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	err = store.IncrementUpload(ip2, 2048, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	// Check counts are separate
	count1, err := store.GetUploadCount(ip1, window)
	if err != nil {
		t.Fatalf("GetUploadCount() error = %v", err)
	}
	if count1 != 1 {
		t.Errorf("GetUploadCount(ip1) = %v, want 1", count1)
	}

	count2, err := store.GetUploadCount(ip2, window)
	if err != nil {
		t.Fatalf("GetUploadCount() error = %v", err)
	}
	if count2 != 1 {
		t.Errorf("GetUploadCount(ip2) = %v, want 1", count2)
	}

	// Check bytes are separate
	bytes1, err := store.GetBytesUsed(ip1, window)
	if err != nil {
		t.Fatalf("GetBytesUsed() error = %v", err)
	}
	if bytes1 != 1024 {
		t.Errorf("GetBytesUsed(ip1) = %v, want 1024", bytes1)
	}

	bytes2, err := store.GetBytesUsed(ip2, window)
	if err != nil {
		t.Fatalf("GetBytesUsed() error = %v", err)
	}
	if bytes2 != 2048 {
		t.Errorf("GetBytesUsed(ip2) = %v, want 2048", bytes2)
	}
}

func TestMemoryStore_CapacityLimit(t *testing.T) {
	// Create store with very small capacity
	store := NewMemoryStore(2, time.Minute)
	defer store.Close()

	window := time.Minute

	// Add uploads up to capacity
	err := store.IncrementUpload("ip1", 1024, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	err = store.IncrementUpload("ip2", 1024, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	// This should fail due to capacity limit
	err = store.IncrementUpload("ip3", 1024, window)
	if err != ErrStoreCapacityExceeded {
		t.Errorf("IncrementUpload() error = %v, want %v", err, ErrStoreCapacityExceeded)
	}
}

func TestMemoryStore_Cleanup(t *testing.T) {
	store := NewMemoryStore(100, time.Minute)
	defer store.Close()

	ip := "203.0.113.1"
	window := 100 * time.Millisecond

	// Add upload
	err := store.IncrementUpload(ip, 1024, window)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	// Wait for expiry
	time.Sleep(150 * time.Millisecond)

	// Run cleanup with the same window as the upload
	if memStore, ok := store.(*memoryStore); ok {
		err = memStore.CleanupWithWindow(window)
		if err != nil {
			t.Fatalf("CleanupWithWindow() error = %v", err)
		}
	} else {
		t.Fatal("Store is not a memory store")
	}

	// Check that expired entries are removed
	count, err := store.GetUploadCount(ip, time.Hour) // Use longer window to see if data was actually removed
	if err != nil {
		t.Fatalf("GetUploadCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetUploadCount() after cleanup = %v, want 0", count)
	}
}

func TestMemoryStore_HealthCheck(t *testing.T) {
	store := NewMemoryStore(100, time.Minute)

	// Should be healthy
	err := store.HealthCheck()
	if err != nil {
		t.Errorf("HealthCheck() error = %v, want nil", err)
	}

	// Close store
	store.Close()

	// Should now be unhealthy
	err = store.HealthCheck()
	if err != ErrStoreClosed {
		t.Errorf("HealthCheck() after close error = %v, want %v", err, ErrStoreClosed)
	}
}

func TestMemoryStore_Close(t *testing.T) {
	store := NewMemoryStore(100, time.Minute)

	// Add some data
	err := store.IncrementUpload("ip1", 1024, time.Minute)
	if err != nil {
		t.Fatalf("IncrementUpload() error = %v", err)
	}

	// Close store
	err = store.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Operations should now fail
	_, err = store.GetUploadCount("ip1", time.Minute)
	if err != ErrStoreClosed {
		t.Errorf("GetUploadCount() after close error = %v, want %v", err, ErrStoreClosed)
	}

	// Double close should be safe
	err = store.Close()
	if err != nil {
		t.Errorf("Double Close() error = %v", err)
	}
}

func BenchmarkMemoryStore_IncrementUpload(b *testing.B) {
	store := NewMemoryStore(10000, time.Minute)
	defer store.Close()

	window := time.Minute

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ip := "203.0.113." + string(rune(i%255))
		_ = store.IncrementUpload(ip, 1024, window)
	}
}

func BenchmarkMemoryStore_GetUploadCount(b *testing.B) {
	store := NewMemoryStore(10000, time.Minute)
	defer store.Close()

	window := time.Minute

	// Pre-populate with data
	for i := 0; i < 1000; i++ {
		ip := "203.0.113." + string(rune(i%255))
		_ = store.IncrementUpload(ip, 1024, window)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ip := "203.0.113." + string(rune(i%255))
		_, _ = store.GetUploadCount(ip, window)
	}
}
