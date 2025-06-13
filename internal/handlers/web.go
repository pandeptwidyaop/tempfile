package handlers

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/tempfile/internal/config"
	"github.com/pandeptwidyaop/tempfile/internal/services"
	"github.com/pandeptwidyaop/tempfile/internal/utils"
)

// WebHandler handles web UI endpoints
type WebHandler struct {
	config          *config.Config
	uploadService   *services.UploadService
	templateService *services.TemplateService
}

// NewWebHandler creates a new web handler instance
func NewWebHandler(cfg *config.Config, uploadSvc *services.UploadService, templateSvc *services.TemplateService) *WebHandler {
	return &WebHandler{
		config:          cfg,
		uploadService:   uploadSvc,
		templateService: templateSvc,
	}
}

// UploadPage renders the upload page
func (h *WebHandler) UploadPage(c *fiber.Ctx) error {
	baseURL := utils.GetBaseURL(c, h.config.PublicURL)
	return h.templateService.RenderUploadPage(c, baseURL)
}

// SuccessPage renders the success page
func (h *WebHandler) SuccessPage(c *fiber.Ctx) error {
	baseURL := utils.GetBaseURL(c, h.config.PublicURL)
	return h.templateService.RenderSuccessPage(c, baseURL)
}

// UploadFileHandler handles both API and web upload
func (h *WebHandler) UploadFileHandler(c *fiber.Ctx) error {
	// Check if this is a web request (has Accept header with text/html)
	acceptHeader := c.Get("Accept")
	isWebRequest := strings.Contains(acceptHeader, "text/html")

	// Use the upload service to process the file
	result, err := h.uploadService.ProcessFileUpload(c)
	if err != nil {
		if isWebRequest {
			return h.templateService.RenderErrorPage(c, "Upload Failed", err.Error(), "Please check your file and try again.")
		}
		return err // Return the fiber error for API
	}

	if isWebRequest {
		// Use url.Values for proper URL encoding
		v := url.Values{}
		v.Set("file", result.Filename)
		v.Set("original", result.OriginalName)
		v.Set("size", fmt.Sprintf("%d", result.Size))
		v.Set("size_human", result.SizeHuman)
		v.Set("expires_at", result.ExpiresAt.Format(time.RFC3339))
		v.Set("expires_in", result.ExpiresIn)

		successURL := "/success?" + v.Encode()
		return c.Redirect(successURL)
	}

	// Return JSON for API requests
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
