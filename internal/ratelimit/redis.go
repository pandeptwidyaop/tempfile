package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisStore implements the Store interface using Redis
type redisStore struct {
	client    *redis.Client
	keyPrefix string
	ctx       context.Context
}

// NewRedisStore creates a new Redis-based rate limit store
func NewRedisStore(redisURL, password string, db int, poolSize, timeout int) (Store, error) {
	// Parse Redis URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis URL: %w", err)
	}

	// Override with provided values
	if password != "" {
		opt.Password = password
	}
	opt.DB = db
	opt.PoolSize = poolSize
	opt.DialTimeout = time.Duration(timeout) * time.Second
	opt.ReadTimeout = time.Duration(timeout) * time.Second
	opt.WriteTimeout = time.Duration(timeout) * time.Second

	client := redis.NewClient(opt)

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis connection failed: %w", err)
	}

	return &redisStore{
		client:    client,
		keyPrefix: "ratelimit:",
		ctx:       ctx,
	}, nil
}

// GetUploadCount returns the number of uploads for an IP within the time window
func (s *redisStore) GetUploadCount(ip string, window time.Duration) (int, error) {
	key := s.keyPrefix + "uploads:" + ip
	cutoff := time.Now().Add(-window).Unix()

	// Remove expired entries and count remaining
	pipe := s.client.Pipeline()
	pipe.ZRemRangeByScore(s.ctx, key, "0", strconv.FormatInt(cutoff, 10))
	pipe.ZCard(s.ctx, key)

	results, err := pipe.Exec(s.ctx)
	if err != nil {
		return 0, fmt.Errorf("Redis pipeline error: %w", err)
	}

	count := results[1].(*redis.IntCmd).Val()
	return int(count), nil
}

// GetBytesUsed returns the total bytes uploaded for an IP within the time window
func (s *redisStore) GetBytesUsed(ip string, window time.Duration) (int64, error) {
	key := s.keyPrefix + "bytes:" + ip
	cutoff := time.Now().Add(-window).Unix()

	// Get all entries within the window
	entries, err := s.client.ZRangeByScoreWithScores(s.ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatInt(cutoff, 10),
		Max: "+inf",
	}).Result()

	if err != nil {
		return 0, fmt.Errorf("Redis range query error: %w", err)
	}

	var totalBytes int64
	for _, entry := range entries {
		// The member is the file size, score is timestamp
		if size, err := strconv.ParseInt(entry.Member.(string), 10, 64); err == nil {
			totalBytes += size
		}
	}

	return totalBytes, nil
}

// IncrementUpload records a new upload for an IP with the given file size
func (s *redisStore) IncrementUpload(ip string, fileSize int64, window time.Duration) error {
	now := time.Now()
	timestamp := now.Unix()

	uploadsKey := s.keyPrefix + "uploads:" + ip
	bytesKey := s.keyPrefix + "bytes:" + ip

	// Use pipeline for atomic operations
	pipe := s.client.Pipeline()

	// Add upload record (score = timestamp, member = timestamp for uniqueness)
	pipe.ZAdd(s.ctx, uploadsKey, redis.Z{
		Score:  float64(timestamp),
		Member: fmt.Sprintf("%d_%d", timestamp, time.Now().UnixNano()),
	})

	// Add bytes record (score = timestamp, member = file size)
	pipe.ZAdd(s.ctx, bytesKey, redis.Z{
		Score:  float64(timestamp),
		Member: strconv.FormatInt(fileSize, 10),
	})

	// Set expiry for keys (window + buffer)
	expiry := window + time.Hour
	pipe.Expire(s.ctx, uploadsKey, expiry)
	pipe.Expire(s.ctx, bytesKey, expiry)

	// Remove old entries
	cutoff := now.Add(-window).Unix()
	pipe.ZRemRangeByScore(s.ctx, uploadsKey, "0", strconv.FormatInt(cutoff, 10))
	pipe.ZRemRangeByScore(s.ctx, bytesKey, "0", strconv.FormatInt(cutoff, 10))

	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return fmt.Errorf("Redis increment error: %w", err)
	}

	return nil
}

// Cleanup removes expired entries from the store
func (s *redisStore) Cleanup() error {
	// Get all rate limit keys
	pattern := s.keyPrefix + "*"
	keys, err := s.client.Keys(s.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("Redis keys scan error: %w", err)
	}

	// Clean up expired entries from each key
	cutoff := time.Now().Add(-24 * time.Hour).Unix() // Clean entries older than 24 hours

	pipe := s.client.Pipeline()
	for _, key := range keys {
		pipe.ZRemRangeByScore(s.ctx, key, "0", strconv.FormatInt(cutoff, 10))
	}

	_, err = pipe.Exec(s.ctx)
	if err != nil {
		return fmt.Errorf("Redis cleanup error: %w", err)
	}

	return nil
}

// HealthCheck verifies the store is functioning properly
func (s *redisStore) HealthCheck() error {
	return s.client.Ping(s.ctx).Err()
}

// Close closes the store and releases resources
func (s *redisStore) Close() error {
	return s.client.Close()
}

// GetStats returns statistics about the Redis store
func (s *redisStore) GetStats() map[string]interface{} {
	info := s.client.Info(s.ctx, "memory", "keyspace").Val()

	// Count rate limit keys
	pattern := s.keyPrefix + "*"
	keys, _ := s.client.Keys(s.ctx, pattern).Result()

	stats := map[string]interface{}{
		"type":        "redis",
		"active_keys": len(keys),
		"redis_info":  info,
	}

	return stats
}

// Lua script for atomic rate limit check and increment
const rateLimitScript = `
local uploads_key = KEYS[1]
local bytes_key = KEYS[2]
local current_time = tonumber(ARGV[1])
local window_seconds = tonumber(ARGV[2])
local file_size = tonumber(ARGV[3])
local upload_limit = tonumber(ARGV[4])
local bytes_limit = tonumber(ARGV[5])

local cutoff = current_time - window_seconds

-- Clean old entries
redis.call('ZREMRANGEBYSCORE', uploads_key, 0, cutoff)
redis.call('ZREMRANGEBYSCORE', bytes_key, 0, cutoff)

-- Get current counts
local upload_count = redis.call('ZCARD', uploads_key)
local bytes_entries = redis.call('ZRANGEBYSCORE', bytes_key, cutoff, '+inf')

local total_bytes = 0
for i = 1, #bytes_entries do
    total_bytes = total_bytes + tonumber(bytes_entries[i])
end

-- Check limits
if upload_count >= upload_limit then
    return {0, upload_count, total_bytes, "upload_limit"}
end

if total_bytes + file_size > bytes_limit then
    return {0, upload_count, total_bytes, "bytes_limit"}
end

-- Add new entries
local unique_id = current_time .. '_' .. redis.call('INCR', 'ratelimit:counter')
redis.call('ZADD', uploads_key, current_time, unique_id)
redis.call('ZADD', bytes_key, current_time, file_size)

-- Set expiry
local expiry = window_seconds + 3600 -- window + 1 hour buffer
redis.call('EXPIRE', uploads_key, expiry)
redis.call('EXPIRE', bytes_key, expiry)

return {1, upload_count + 1, total_bytes + file_size, "ok"}
`

// AtomicCheckAndIncrement performs atomic rate limit check and increment
func (s *redisStore) AtomicCheckAndIncrement(ip string, fileSize int64, window time.Duration, uploadLimit int, bytesLimit int64) (bool, int, int64, string, error) {
	uploadsKey := s.keyPrefix + "uploads:" + ip
	bytesKey := s.keyPrefix + "bytes:" + ip

	result, err := s.client.Eval(s.ctx, rateLimitScript, []string{uploadsKey, bytesKey},
		time.Now().Unix(),
		int64(window.Seconds()),
		fileSize,
		uploadLimit,
		bytesLimit,
	).Result()

	if err != nil {
		return false, 0, 0, "", fmt.Errorf("Redis atomic operation error: %w", err)
	}

	results := result.([]interface{})
	allowed := results[0].(int64) == 1
	uploadCount := int(results[1].(int64))
	totalBytes := results[2].(int64)
	reason := results[3].(string)

	return allowed, uploadCount, totalBytes, reason, nil
}
