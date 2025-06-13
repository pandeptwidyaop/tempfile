// internal/services/static.go - Hybrid embedded + filesystem loading
package services

import (
	"embed"
	"io/fs"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/tempfile/internal/config"
)

// Embed static files at compile time
// go:embed ../../web/static/*
var staticFiles embed.FS

// StaticService handles serving static files with hybrid loading
type StaticService struct {
	staticFS        fs.FS
	useFileSystem   bool
	staticDir       string
	embeddedFS      fs.FS
}

// NewStaticService creates a new static service with hybrid loading
func NewStaticService(cfg *config.Config) (*StaticService, error) {
	// Determine if we should use filesystem or embedded assets
	useFileSystem := cfg.Debug || os.Getenv("USE_FILESYSTEM_ASSETS") == "true"
	
	// Setup embedded filesystem
	embeddedFS, err := fs.Sub(staticFiles, "web/static")
	if err != nil {
		log.Printf("Warning: could not create static subdirectory from embedded, using full embedded FS: %v", err)
		embeddedFS = staticFiles
	}

	service := &StaticService{
		useFileSystem: useFileSystem,
		staticDir:     cfg.StaticDir,
		embeddedFS:    embeddedFS,
	}

	if useFileSystem {
		// Check if static directory exists
		if _, err := os.Stat(cfg.StaticDir); os.IsNotExist(err) {
			log.Printf("⚠️  Static directory %s not found, falling back to embedded assets", cfg.StaticDir)
			service.useFileSystem = false
		} else {
			log.Printf("🔧 Development mode: Loading static files from %s", cfg.StaticDir)
		}
	}

	if !service.useFileSystem {
		log.Println("✅ Production mode: Using embedded static assets")
	}

	service.staticFS = embeddedFS
	return service, nil
}

// Handler returns a Fiber handler for serving static files
func (s *StaticService) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return s.ServeFile(c)
	}
}

// ServeFile serves a static file with hybrid loading (filesystem first, then embedded)
func (s *StaticService) ServeFile(c *fiber.Ctx) error {
	// Get the file path from URL parameters
	filePath := c.Params("*")
	if filePath == "" {
		return c.Status(404).SendString("File not found")
	}

	// Clean the path to prevent directory traversal
	filePath = filepath.Clean(filePath)
	if filePath == "." || filePath == "/" {
		return c.Status(404).SendString("File not found")
	}

	// Remove leading slash if present
	filePath = strings.TrimPrefix(filePath, "/")

	var fileData []byte
	var err error
	var source string

	// Try filesystem first if enabled
	if s.useFileSystem {
		fullPath := filepath.Join(s.staticDir, filePath)
		fileData, err = os.ReadFile(fullPath)
		if err == nil {
			source = "filesystem"
			log.Printf("🔧 Served from filesystem: %s", filePath)
		} else {
			// Fallback to embedded if filesystem fails
			log.Printf("📁 File not found in filesystem, trying embedded: %s", filePath)
		}
	}

	// If filesystem failed or not enabled, try embedded
	if fileData == nil {
		fileData, err = fs.ReadFile(s.embeddedFS, filePath)
		if err != nil {
			log.Printf("❌ Static file not found in both filesystem and embedded: %s", filePath)
			return c.Status(404).SendString("File not found")
		}
		source = "embedded"
	}

	// Determine content type based on file extension
	contentType := s.getContentType(filePath)

	// Set appropriate headers
	c.Set("Content-Type", contentType)
	
	// Different cache headers based on source
	if source == "filesystem" {
		// Shorter cache for development
		c.Set("Cache-Control", "no-cache, must-revalidate")
	} else {
		// Longer cache for embedded assets
		c.Set("Cache-Control", "public, max-age=3600")
	}
	
	// Add security headers for assets
	if strings.HasSuffix(filePath, ".js") {
		c.Set("X-Content-Type-Options", "nosniff")
	}

	if source == "embedded" || !s.useFileSystem {
		log.Printf("✅ Served from %s: %s (%d bytes, %s)", source, filePath, len(fileData), contentType)
	}
	
	return c.Send(fileData)
}

// getContentType determines content type based on file extension
func (s *StaticService) getContentType(filePath string) string {
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		// Default content type based on file extension
		switch filepath.Ext(filePath) {
		case ".css":
			contentType = "text/css"
		case ".js":
			contentType = "application/javascript"
		case ".html":
			contentType = "text/html"
		case ".json":
			contentType = "application/json"
		case ".svg":
			contentType = "image/svg+xml"
		case ".png":
			contentType = "image/png"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".gif":
			contentType = "image/gif"
		case ".ico":
			contentType = "image/x-icon"
		case ".woff", ".woff2":
			contentType = "font/woff"
		case ".ttf":
			contentType = "font/ttf"
		case ".eot":
			contentType = "application/vnd.ms-fontobject"
		default:
			contentType = "application/octet-stream"
		}
	}
	return contentType
}

// GetMode returns the current loading mode
func (s *StaticService) GetMode() string {
	if s.useFileSystem {
		return "filesystem"
	}
	return "embedded"
}

// ListFiles returns a list of available static files (for debugging)
func (s *StaticService) ListFiles() ([]string, error) {
	var files []string
	
	if s.useFileSystem {
		// List from filesystem
		err := filepath.Walk(s.staticDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				relPath, _ := filepath.Rel(s.staticDir, path)
				files = append(files, relPath)
			}
			return nil
		})
		return files, err
	} else {
		// List from embedded
		err := fs.WalkDir(s.embeddedFS, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		return files, err
	}
}
