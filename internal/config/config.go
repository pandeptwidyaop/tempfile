package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server config
	Port      string
	PublicURL string
	
	// File storage config
	UploadDir       string
	MaxFileSize     int64
	FileExpiryHours int
	
	// Middleware config
	EnableCORS    bool
	CORSOrigins   string
	EnableLogging bool
	
	// Cleanup config
	CleanupIntervalSeconds int
	
	// Environment config
	AppEnv   string
	Debug    bool
	LogLevel string
	
	// Web UI config
	EnableWebUI   bool
	StaticDir     string
	TemplatesDir  string
	DefaultTheme  string
	
	// Rate limiting config
	EnableRateLimit              bool
	RateLimitStore              string
	RateLimitUploadsPerMinute   int
	RateLimitBytesPerHour       int64
	RateLimitWindowMinutes      int
	RateLimitTrustedProxies     []string
	RateLimitIPHeaders          []string
	RateLimitWhitelistIPs       []string
	RateLimitCustomLimits       map[string]RateLimitEndpointConfig
	
	// Redis config
	RedisURL                    string
	RedisPassword               string
	RedisDB                     int
	RedisPoolSize               int
	RedisTimeout                int
}

// RateLimitEndpointConfig holds custom rate limits for specific endpoints
type RateLimitEndpointConfig struct {
	UploadsPerMinute int
	BytesPerHour     int64
	WindowMinutes    int
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables and defaults")
	}

	config := &Config{
		// Server config
		Port:      getEnvOrDefault("PORT", "3000"),
		PublicURL: getEnvOrDefault("PUBLIC_URL", "http://localhost:3000"),
		
		// File storage config
		UploadDir:       getEnvOrDefault("UPLOAD_DIR", "./uploads"),
		MaxFileSize:     getEnvAsInt64OrDefault("MAX_FILE_SIZE", 100*1024*1024), // 100MB
		FileExpiryHours: getEnvAsIntOrDefault("FILE_EXPIRY_HOURS", 1),
		
		// Middleware config
		EnableCORS:    getEnvAsBoolOrDefault("ENABLE_CORS", true),
		CORSOrigins:   getEnvOrDefault("CORS_ORIGINS", "*"),
		EnableLogging: getEnvAsBoolOrDefault("ENABLE_LOGGING", true),
		
		// Cleanup config
		CleanupIntervalSeconds: getEnvAsIntOrDefault("CLEANUP_INTERVAL_SECONDS", 1),
		
		// Environment config
		AppEnv:   getEnvOrDefault("APP_ENV", "production"),
		Debug:    getEnvAsBoolOrDefault("DEBUG", false),
		LogLevel: getEnvOrDefault("LOG_LEVEL", "info"),
		
		// Web UI config
		EnableWebUI:  getEnvAsBoolOrDefault("ENABLE_WEB_UI", true),
		StaticDir:    getEnvOrDefault("STATIC_DIR", "./web/static"),
		TemplatesDir: getEnvOrDefault("TEMPLATES_DIR", "./web/templates"),
		DefaultTheme: getEnvOrDefault("DEFAULT_THEME", "dark"),
		
		// Rate limiting config
		EnableRateLimit:              getEnvAsBoolOrDefault("ENABLE_RATE_LIMIT", false),
		RateLimitStore:              getEnvOrDefault("RATE_LIMIT_STORE", "memory"),
		RateLimitUploadsPerMinute:   getEnvAsIntOrDefault("RATE_LIMIT_UPLOADS_PER_MINUTE", 5),
		RateLimitBytesPerHour:       getEnvAsInt64OrDefault("RATE_LIMIT_BYTES_PER_HOUR", 100*1024*1024), // 100MB
		RateLimitWindowMinutes:      getEnvAsIntOrDefault("RATE_LIMIT_WINDOW_MINUTES", 60),
		RateLimitTrustedProxies:     getEnvAsStringSliceOrDefault("RATE_LIMIT_TRUSTED_PROXIES", []string{"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}),
		RateLimitIPHeaders:          getEnvAsStringSliceOrDefault("RATE_LIMIT_IP_HEADERS", []string{"CF-Connecting-IP", "X-Real-IP", "X-Forwarded-For"}),
		RateLimitWhitelistIPs:       getEnvAsStringSliceOrDefault("RATE_LIMIT_WHITELIST_IPS", []string{}),
		RateLimitCustomLimits:       parseCustomRateLimits(),
		
		// Redis config
		RedisURL:      getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsIntOrDefault("REDIS_DB", 0),
		RedisPoolSize: getEnvAsIntOrDefault("REDIS_POOL_SIZE", 10),
		RedisTimeout:  getEnvAsIntOrDefault("REDIS_TIMEOUT", 5),
	}

	// Add colon prefix to port if not present
	if !strings.HasPrefix(config.Port, ":") {
		config.Port = ":" + config.Port
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}

// Helper functions for environment variable parsing

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64OrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsStringSliceOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// ValidateTrustedProxies validates the trusted proxy configuration
func (c *Config) ValidateTrustedProxies() error {
	for _, proxy := range c.RateLimitTrustedProxies {
		// Try to parse as IP first
		if ip := net.ParseIP(proxy); ip != nil {
			continue
		}
		
		// Try to parse as CIDR
		if _, _, err := net.ParseCIDR(proxy); err != nil {
			return fmt.Errorf("invalid trusted proxy format: %s", proxy)
		}
	}
	return nil
}

// parseCustomRateLimits parses custom rate limits from environment variables
func parseCustomRateLimits() map[string]RateLimitEndpointConfig {
	customLimits := make(map[string]RateLimitEndpointConfig)
	
	// Parse format: ENDPOINT_PATH:uploads_per_min:bytes_per_hour:window_min
	// Example: RATE_LIMIT_CUSTOM_ENDPOINTS="/api/upload:10:209715200:30,/bulk:2:52428800:60"
	customEndpoints := getEnvOrDefault("RATE_LIMIT_CUSTOM_ENDPOINTS", "")
	if customEndpoints == "" {
		return customLimits
	}
	
	endpoints := strings.Split(customEndpoints, ",")
	for _, endpoint := range endpoints {
		parts := strings.Split(strings.TrimSpace(endpoint), ":")
		if len(parts) != 4 {
			continue
		}
		
		path := parts[0]
		uploadsPerMin, err1 := strconv.Atoi(parts[1])
		bytesPerHour, err2 := strconv.ParseInt(parts[2], 10, 64)
		windowMin, err3 := strconv.Atoi(parts[3])
		
		if err1 == nil && err2 == nil && err3 == nil {
			customLimits[path] = RateLimitEndpointConfig{
				UploadsPerMinute: uploadsPerMin,
				BytesPerHour:     bytesPerHour,
				WindowMinutes:    windowMin,
			}
		}
	}
	
	return customLimits
}
