package handlers

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/tempfile/internal/config"
	"github.com/pandeptwidyaop/tempfile/internal/utils"
)

// FileHandler handles file operations
type FileHandler struct {
	config *config.Config
}

// NewFileHandler creates a new file handler instance
func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{
		config: cfg,
	}
}

// DownloadFile handles file download
func (h *FileHandler) DownloadFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filePath := filepath.Join(h.config.UploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Check if file has expired using utility function
	expired, err := utils.IsFileExpired(filename, time.Now())
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid filename format",
		})
	}

	if expired {
		// Remove expired file
		_ = os.Remove(filePath)
		return c.Status(404).JSON(fiber.Map{
			"error": "File has expired",
		})
	}

	if h.config.Debug {
		log.Printf("File downloaded: %s", filename)
	}

	// Download file
	return c.SendFile(filePath)
}
