package ratelimit

import (
	"testing"
	"time"
)

// Note: These tests require a running Redis instance
// Skip if Redis is not available

func TestRedisStore_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis integration test in short mode")
	}
	
	// Try to create Redis store
	store, err := NewRedisStore("redis://localhost:6379", "", 0, 10, 5)
	if err != nil {
		t.Skipf("Redis not available, skipping test: %v", err)
	}
	defer store.Close()
	
	ip := "203.0.113.1"
	window := time.Minute
	
	// Test basic operations
	t.Run("GetUploadCount_Initially_Zero", func(t *testing.T) {
		count, err := store.GetUploadCount(ip, window)
		if err != nil {
			t.Fatalf("GetUploadCount() error = %v", err)
		}
		if count != 0 {
			t.Errorf("GetUploadCount() = %v, want 0", count)
		}
	})
	
	t.Run("IncrementUpload_And_GetCount", func(t *testing.T) {
		err := store.IncrementUpload(ip, 1024, window)
		if err != nil {
			t.Fatalf("IncrementUpload() error = %v", err)
		}
		
		count, err := store.GetUploadCount(ip, window)
		if err != nil {
			t.Fatalf("GetUploadCount() error = %v", err)
		}
		if count != 1 {
			t.Errorf("GetUploadCount() = %v, want 1", count)
		}
	})
	
	t.Run("GetBytesUsed", func(t *testing.T) {
		bytes, err := store.GetBytesUsed(ip, window)
		if err != nil {
			t.Fatalf("GetBytesUsed() error = %v", err)
		}
		if bytes != 1024 {
			t.Errorf("GetBytesUsed() = %v, want 1024", bytes)
		}
	})
	
	t.Run("HealthCheck", func(t *testing.T) {
		err := store.HealthCheck()
		if err != nil {
			t.Errorf("HealthCheck() error = %v", err)
		}
	})
}

func TestRedisStore_AtomicOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis integration test in short mode")
	}
	
	store, err := NewRedisStore("redis://localhost:6379", "", 0, 10, 5)
	if err != nil {
		t.Skipf("Redis not available, skipping test: %v", err)
	}
	defer store.Close()
	
	redisStore := store.(*redisStore)
	ip := "203.0.113.2"
	window := time.Minute
	
	t.Run("AtomicCheckAndIncrement_WithinLimits", func(t *testing.T) {
		allowed, uploadCount, totalBytes, reason, err := redisStore.AtomicCheckAndIncrement(
			ip, 1024, window, 5, 10240)
		
		if err != nil {
			t.Fatalf("AtomicCheckAndIncrement() error = %v", err)
		}
		
		if !allowed {
			t.Errorf("AtomicCheckAndIncrement() allowed = %v, want true", allowed)
		}
		
		if uploadCount != 1 {
			t.Errorf("AtomicCheckAndIncrement() uploadCount = %v, want 1", uploadCount)
		}
		
		if totalBytes != 1024 {
			t.Errorf("AtomicCheckAndIncrement() totalBytes = %v, want 1024", totalBytes)
		}
		
		if reason != "ok" {
			t.Errorf("AtomicCheckAndIncrement() reason = %v, want 'ok'", reason)
		}
	})
	
	t.Run("AtomicCheckAndIncrement_ExceedsUploadLimit", func(t *testing.T) {
		// Add more uploads to exceed limit
		for i := 0; i < 5; i++ {
			redisStore.AtomicCheckAndIncrement(ip, 100, window, 5, 10240)
		}
		
		allowed, _, _, reason, err := redisStore.AtomicCheckAndIncrement(
			ip, 100, window, 5, 10240)
		
		if err != nil {
			t.Fatalf("AtomicCheckAndIncrement() error = %v", err)
		}
		
		if allowed {
			t.Errorf("AtomicCheckAndIncrement() allowed = %v, want false", allowed)
		}
		
		if reason != "upload_limit" {
			t.Errorf("AtomicCheckAndIncrement() reason = %v, want 'upload_limit'", reason)
		}
	})
}

func BenchmarkRedisStore_AtomicCheckAndIncrement(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping Redis benchmark in short mode")
	}
	
	store, err := NewRedisStore("redis://localhost:6379", "", 0, 10, 5)
	if err != nil {
		b.Skipf("Redis not available, skipping benchmark: %v", err)
	}
	defer store.Close()
	
	redisStore := store.(*redisStore)
	window := time.Minute
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ip := "bench.test." + string(rune(i%1000))
		redisStore.AtomicCheckAndIncrement(ip, 1024, window, 100, 1048576)
	}
}