package services

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/tempfile/internal/config"
	"github.com/pandeptwidyaop/tempfile/internal/models"
)

// Embed templates at compile time
// go:embed ../../web/templates/*.html
var templateFiles embed.FS

// TemplateService handles HTML template rendering with hybrid loading
type TemplateService struct {
	config          *config.Config
	templates       *template.Template
	useFileSystem   bool
}

// NewTemplateService creates a new template service instance with hybrid loading
func NewTemplateService(cfg *config.Config) *TemplateService {
	return &TemplateService{
		config: cfg,
	}
}

// Initialize initializes HTML templates with hybrid loading (filesystem first, then embedded)
func (s *TemplateService) Initialize() error {
	// Determine if we should use filesystem or embedded templates
	s.useFileSystem = s.config.Debug || os.Getenv("USE_FILESYSTEM_ASSETS") == "true"

	var tmpl *template.Template
	var err error

	// Try filesystem first if enabled
	if s.useFileSystem {
		if _, statErr := os.Stat(s.config.TemplatesDir); os.IsNotExist(statErr) {
			log.Printf("⚠️  Templates directory %s not found, falling back to embedded templates", s.config.TemplatesDir)
			s.useFileSystem = false
		} else {
			// Load from filesystem
			templatePattern := filepath.Join(s.config.TemplatesDir, "*.html")
			tmpl, err = template.ParseGlob(templatePattern)
			if err != nil {
				log.Printf("⚠️  Failed to load templates from filesystem: %v, falling back to embedded", err)
				s.useFileSystem = false
			} else {
				log.Printf("🔧 Development mode: Loading templates from %s", s.config.TemplatesDir)
			}
		}
	}

	// If filesystem failed or not enabled, use embedded templates
	if !s.useFileSystem || tmpl == nil {
		// Try to get templates subdirectory from embedded files
		templatesFS, err := fs.Sub(templateFiles, "web/templates")
		if err != nil {
			// If sub fails, use the full embedded FS and parse with full paths
			tmpl, err = template.ParseFS(templateFiles, "web/templates/*.html")
			if err != nil {
				return fmt.Errorf("failed to parse embedded templates: %w", err)
			}
			log.Println("✅ Production mode: Templates loaded from embedded files (full path)")
		} else {
			// Parse templates from subdirectory
			tmpl, err = template.ParseFS(templatesFS, "*.html")
			if err != nil {
				return fmt.Errorf("failed to parse embedded templates from subdirectory: %w", err)
			}
			log.Println("✅ Production mode: Templates loaded from embedded files")
		}
	}

	s.templates = tmpl
	return nil
}

// Render renders an HTML template
func (s *TemplateService) Render(c *fiber.Ctx, templateName string, data interface{}) error {
	c.Set("Content-Type", "text/html")

	// In development mode with filesystem templates, reload templates on each request
	if s.useFileSystem && s.config.Debug {
		if err := s.reloadTemplatesIfNeeded(); err != nil {
			log.Printf("Warning: Failed to reload templates: %v", err)
		}
	}

	// Create a buffer to render the template
	var buf strings.Builder
	err := s.templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		log.Printf("Template rendering error: %v", err)
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.SendString(buf.String())
}

// reloadTemplatesIfNeeded reloads templates from filesystem in development mode
func (s *TemplateService) reloadTemplatesIfNeeded() error {
	if !s.useFileSystem {
		return nil
	}

	templatePattern := filepath.Join(s.config.TemplatesDir, "*.html")
	tmpl, err := template.ParseGlob(templatePattern)
	if err != nil {
		return err
	}

	s.templates = tmpl
	return nil
}

// GetMode returns the current loading mode
func (s *TemplateService) GetMode() string {
	if s.useFileSystem {
		return "filesystem"
	}
	return "embedded"
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
