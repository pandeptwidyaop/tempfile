package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/tempfile/internal/config"
	"github.com/pandeptwidyaop/tempfile/internal/models"
	"github.com/pandeptwidyaop/tempfile/internal/services"
)

// APIHandler handles API endpoints
type APIHandler struct {
	config        *config.Config
	uploadService *services.UploadService
}

// NewAPIHandler creates a new API handler instance
func NewAPIHandler(cfg *config.Config, uploadSvc *services.UploadService) *APIHandler {
	return &APIHandler{
		config:        cfg,
		uploadService: uploadSvc,
	}
}

// UploadFile handles API file upload
func (h *APIHandler) UploadFile(c *fiber.Ctx) error {
	result, err := h.uploadService.ProcessFileUpload(c)
	if err != nil {
		return err
	}

	// Convert to API response format
	response := fiber.Map{
		"message":       result.Message,
		"filename":      result.Filename,
		"original_name": result.OriginalName,
		"size":          result.Size,
		"size_human":    result.SizeHuman,
		"expires_at":    result.ExpiresAt.Format(time.RFC3339),
		"expires_in":    result.ExpiresIn,
		"download_url":  result.DownloadURL,
	}

	return c.JSON(response)
}

// HealthCheck handles health check endpoint
func (h *APIHandler) HealthCheck(c *fiber.Ctx) error {
	response := models.HealthResponse{
		Status:      "healthy",
		Environment: h.config.AppEnv,
		Version:     "1.0.0",
		Uptime:      time.Since(startTime).String(),
	}

	return c.JSON(response)
}

var startTime = time.Now()
