package services

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/tempfile/internal/config"
	"github.com/pandeptwidyaop/tempfile/internal/models"
	"github.com/pandeptwidyaop/tempfile/internal/utils"
)

// UploadService handles file upload operations
type UploadService struct {
	config *config.Config
}

// NewUploadService creates a new upload service instance
func NewUploadService(cfg *config.Config) *UploadService {
	return &UploadService{
		config: cfg,
	}
}

// ProcessFileUpload handles the core file upload logic
func (s *UploadService) ProcessFileUpload(c *fiber.Ctx) (*models.UploadResponse, error) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return nil, fiber.NewError(400, "No file uploaded")
	}

	// Check file size
	if file.Size > s.config.MaxFileSize {
		return nil, fiber.NewError(400, fmt.Sprintf("File size exceeds %s limit", utils.FormatBytes(s.config.MaxFileSize)))
	}

	// Generate filename based on unix timestamp (now + expiry hours) + extension
	expiryTime := time.Now().Add(time.Duration(s.config.FileExpiryHours) * time.Hour)
	filename := utils.GenerateFilename(file.Filename, expiryTime)

	// Save file with unix timestamp + extension as filename
	filePath := filepath.Join(s.config.UploadDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		log.Printf("Error saving file: %v", err)
		return nil, fiber.NewError(500, "Failed to save file")
	}

	if s.config.Debug {
		log.Printf("File uploaded: %s (original: %s, size: %s)", filename, file.Filename, utils.FormatBytes(file.Size))
	}

	response := &models.UploadResponse{
		Message:      "File uploaded successfully",
		Filename:     filename,
		OriginalName: file.Filename,
		Size:         file.Size,
		SizeHuman:    utils.FormatBytes(file.Size),
		ExpiresAt:    expiryTime,
		ExpiresIn:    fmt.Sprintf("%d hour(s)", s.config.FileExpiryHours),
		DownloadURL:  fmt.Sprintf("/%s", filename),
	}

	return response, nil
}
