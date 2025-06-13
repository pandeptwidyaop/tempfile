package services

import (
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/tempfile/internal/config"
	"github.com/pandeptwidyaop/tempfile/internal/models"
)

// TemplateService handles HTML template rendering
type TemplateService struct {
	config    *config.Config
	templates *template.Template
}

// NewTemplateService creates a new template service instance
func NewTemplateService(cfg *config.Config) *TemplateService {
	return &TemplateService{
		config: cfg,
	}
}

// Initialize initializes HTML templates
func (s *TemplateService) Initialize() error {
	templatePattern := filepath.Join(s.config.TemplatesDir, "*.html")
	tmpl, err := template.ParseGlob(templatePattern)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}
	s.templates = tmpl
	return nil
}

// Render renders an HTML template
func (s *TemplateService) Render(c *fiber.Ctx, templateName string, data interface{}) error {
	c.Set("Content-Type", "text/html")

	// Create a buffer to render the template
	var buf strings.Builder
	err := s.templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		log.Printf("Template rendering error: %v", err)
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.SendString(buf.String())
}

// RenderUploadPage renders the upload page
func (s *TemplateService) RenderUploadPage(c *fiber.Ctx, baseURL string) error {
	data := models.WebPageData{
		Title:             "Upload File",
		Theme:             s.config.DefaultTheme,
		FileExpiryHours:   s.config.FileExpiryHours,
		MaxFileSizeHuman:  s.formatBytes(s.config.MaxFileSize),
		BaseURL:           baseURL,
	}

	return s.Render(c, "upload.html", data)
}

// RenderSuccessPage renders the success page
func (s *TemplateService) RenderSuccessPage(c *fiber.Ctx, baseURL string) error {
	// Get upload result from query params
	filename := c.Query("file")
	originalName := c.Query("original")
	size := c.Query("size")
	sizeHuman := c.Query("size_human")
	expiresAt := c.Query("expires_at")
	expiresIn := c.Query("expires_in")

	if filename == "" {
		// Redirect to home if no file info
		return c.Redirect("/")
	}

	data := models.WebPageData{
		Title:        "Upload Successful",
		Theme:        s.config.DefaultTheme,
		Filename:     filename,
		OriginalName: originalName,
		Size:         size,
		SizeHuman:    sizeHuman,
		ExpiresAt:    expiresAt,
		ExpiresIn:    expiresIn,
		BaseURL:      baseURL,
	}

	return s.Render(c, "success.html", data)
}

// RenderErrorPage renders an error page
func (s *TemplateService) RenderErrorPage(c *fiber.Ctx, title, message, detail string) error {
	data := models.WebPageData{
		Title:        "Error",
		Theme:        s.config.DefaultTheme,
		ErrorTitle:   title,
		ErrorMessage: message,
		ErrorDetail:  detail,
	}

	c.Status(400)
	return s.Render(c, "error.html", data)
}

// Helper method to format bytes (avoiding import cycle)
func (s *TemplateService) formatBytes(bytes int64) string {
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
