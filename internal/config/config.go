package config

import (
	"log"
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
