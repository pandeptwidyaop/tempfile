package services

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pandeptwidyaop/tempfile/internal/config"
	"github.com/pandeptwidyaop/tempfile/internal/utils"
)

// CleanupService handles expired file cleanup
type CleanupService struct {
	config *config.Config
}

// NewCleanupService creates a new cleanup service instance
func NewCleanupService(cfg *config.Config) *CleanupService {
	return &CleanupService{
		config: cfg,
	}
}

// Start starts the cleanup routine
func (s *CleanupService) Start() {
	ticker := time.NewTicker(time.Duration(s.config.CleanupIntervalSeconds) * time.Second)
	defer ticker.Stop()

	log.Printf("🧹 Cleanup routine started (interval: %d second(s))", s.config.CleanupIntervalSeconds)

	for range ticker.C {
		s.cleanupExpiredFiles()
	}
}

// cleanupExpiredFiles removes expired files from the upload directory
func (s *CleanupService) cleanupExpiredFiles() {
	files, err := os.ReadDir(s.config.UploadDir)
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

		filename := file.Name()

		// Check if file has expired using utility function
		expired, err := utils.IsFileExpired(filename, now)
		if err != nil {
			// Skip files that are not in valid timestamp format
			continue
		}

		if expired {
			filePath := filepath.Join(s.config.UploadDir, filename)
			if err := os.Remove(filePath); err != nil {
				log.Printf("Error removing expired file %s: %v", filename, err)
			} else {
				cleanedCount++
				if s.config.Debug {
					log.Printf("Removed expired file: %s", filename)
				}
			}
		}
	}

	if cleanedCount > 0 {
		log.Printf("🗑️  Cleaned up %d expired file(s)", cleanedCount)
	}
}
