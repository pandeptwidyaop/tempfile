package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Port                   string
	UploadDir              string
	MaxFileSize            int64
	FileExpiryHours        int
	EnableCORS             bool
	CORSOrigins            string
	EnableLogging          bool
	CleanupIntervalSeconds int
	AppEnv                 string
	Debug                  bool
	LogLevel               string
}

// loadConfig loads configuration from environment variables and .env file
func loadConfig() *Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables and defaults")
	}

	config := &Config{
		Port:                   getEnvOrDefault("PORT", "3000"),
		UploadDir:              getEnvOrDefault("UPLOAD_DIR", "./uploads"),
		MaxFileSize:            getEnvAsInt64OrDefault("MAX_FILE_SIZE", 100*1024*1024), // 100MB
		FileExpiryHours:        getEnvAsIntOrDefault("FILE_EXPIRY_HOURS", 1),
		EnableCORS:             getEnvAsBoolOrDefault("ENABLE_CORS", true),
		CORSOrigins:            getEnvOrDefault("CORS_ORIGINS", "*"),
		EnableLogging:          getEnvAsBoolOrDefault("ENABLE_LOGGING", true),
		CleanupIntervalSeconds: getEnvAsIntOrDefault("CLEANUP_INTERVAL_SECONDS", 1),
		AppEnv:                 getEnvOrDefault("APP_ENV", "production"),
		Debug:                  getEnvAsBoolOrDefault("DEBUG", false),
		LogLevel:               getEnvOrDefault("LOG_LEVEL", "info"),
	}

	// Add colon prefix to port if not present
	if !strings.HasPrefix(config.Port, ":") {
		config.Port = ":" + config.Port
	}

	return config
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

// Global config variable
var config *Config

func main() {
	// Load configuration
	config = loadConfig()

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: int(config.MaxFileSize),
	})

	// Conditionally add middleware based on config
	if config.EnableLogging {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
		}))
	}

	if config.EnableCORS {
		corsConfig := cors.New()
		if config.CORSOrigins != "*" {
			corsConfig = cors.New(cors.Config{
				AllowOrigins: config.CORSOrigins,
			})
		}
		app.Use(corsConfig)
	}

	// Routes
	app.Post("/", uploadFile)
	app.Get("/:filename", downloadFile)

	// Add health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "healthy",
			"environment": config.AppEnv,
			"version":     "1.0.0",
			"uptime":      time.Since(startTime).String(),
		})
	})

	// Start cleanup routine
	go startCleanupRoutine()

	// Print startup information
	printStartupInfo()

	// Start server
	log.Printf("🚀 Server starting on %s", config.Port)
	log.Fatal(app.Listen(config.Port))
}

var startTime = time.Now()

// printStartupInfo prints configuration information at startup
func printStartupInfo() {
	log.Println("📁 TempFiles Server Configuration:")
	log.Printf("   Environment: %s", config.AppEnv)
	log.Printf("   Port: %s", config.Port)
	log.Printf("   Upload Directory: %s", config.UploadDir)
	log.Printf("   Max File Size: %s", formatBytes(config.MaxFileSize))
	log.Printf("   File Expiry: %d hour(s)", config.FileExpiryHours)
	log.Printf("   Cleanup Interval: %d second(s)", config.CleanupIntervalSeconds)
	log.Printf("   CORS Enabled: %v", config.EnableCORS)
	log.Printf("   Logging Enabled: %v", config.EnableLogging)
	log.Printf("   Debug Mode: %v", config.Debug)
}

// formatBytes converts bytes to human readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Handler for file upload
func uploadFile(c *fiber.Ctx) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Check file size
	if file.Size > config.MaxFileSize {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("File size exceeds %s limit", formatBytes(config.MaxFileSize)),
		})
	}

	// Generate filename based on unix timestamp (now + expiry hours) + extension
	expiryTime := time.Now().Add(time.Duration(config.FileExpiryHours) * time.Hour)

	// Get extension from original file
	originalExt := filepath.Ext(file.Filename)
	if originalExt == "" {
		originalExt = ".bin" // default extension if none
	}

	filename := strconv.FormatInt(expiryTime.Unix(), 10) + originalExt

	// Save file with unix timestamp + extension as filename
	filePath := filepath.Join(config.UploadDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		log.Printf("Error saving file: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	if config.Debug {
		log.Printf("File uploaded: %s (original: %s, size: %s)", filename, file.Filename, formatBytes(file.Size))
	}

	return c.JSON(fiber.Map{
		"message":       "File uploaded successfully",
		"filename":      filename,
		"original_name": file.Filename,
		"size":          file.Size,
		"size_human":    formatBytes(file.Size),
		"expires_at":    expiryTime.Format(time.RFC3339),
		"expires_in":    fmt.Sprintf("%d hour(s)", config.FileExpiryHours),
		"download_url":  fmt.Sprintf("/%s", filename),
	})
}

// Handler for file download
func downloadFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filePath := filepath.Join(config.UploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Parse unix timestamp from filename for expiry validation
	// Remove extension to get unix timestamp
	filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	unixTime, err := strconv.ParseInt(filenameWithoutExt, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid filename format",
		})
	}

	// Check if file has expired
	expiryTime := time.Unix(unixTime, 0)
	if time.Now().After(expiryTime) {
		// Remove expired file
		os.Remove(filePath)
		return c.Status(404).JSON(fiber.Map{
			"error": "File has expired",
		})
	}

	if config.Debug {
		log.Printf("File downloaded: %s", filename)
	}

	// Download file
	return c.SendFile(filePath)
}

// Goroutine to clean up expired files
func startCleanupRoutine() {
	ticker := time.NewTicker(time.Duration(config.CleanupIntervalSeconds) * time.Second)
	defer ticker.Stop()

	log.Printf("🧹 Cleanup routine started (interval: %d second(s))", config.CleanupIntervalSeconds)

	for {
		select {
		case <-ticker.C:
			cleanupExpiredFiles()
		}
	}
}

// Function to clean up expired files
func cleanupExpiredFiles() {
	files, err := os.ReadDir(config.UploadDir)
	if err != nil {
		log.Printf("Error reading upload directory: %v", err)
		return
	}

	now := time.Now()
	cleanedCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Parse unix timestamp from filename
		filename := file.Name()
		filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
		unixTime, err := strconv.ParseInt(filenameWithoutExt, 10, 64)
		if err != nil {
			// Skip files that are not in unix timestamp format
			continue
		}

		// Check if file has expired
		expiryTime := time.Unix(unixTime, 0)
		if now.After(expiryTime) {
			filePath := filepath.Join(config.UploadDir, filename)
			if err := os.Remove(filePath); err != nil {
				log.Printf("Error removing expired file %s: %v", filename, err)
			} else {
				cleanedCount++
				if config.Debug {
					log.Printf("Removed expired file: %s", filename)
				}
			}
		}
	}

	if cleanedCount > 0 {
		log.Printf("🗑️  Cleaned up %d expired file(s)", cleanedCount)
	}
}
