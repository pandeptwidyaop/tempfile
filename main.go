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
)

const (
	uploadDir   = "./uploads"
	maxFileSize = 100 * 1024 * 1024 // 100MB
	fileExpiry  = 1 * time.Hour     // 1 jam
)

// getPort returns the port from environment variable or default
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return ":" + port
}

func main() {
	// Buat direktori uploads jika belum ada
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	// Inisialisasi Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: maxFileSize,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Post("/", uploadFile)
	app.Get("/:filename", downloadFile)

	// Mulai goroutine untuk membersihkan file expired
	go startCleanupRoutine()

	// Start server
	port := getPort()
	log.Printf("Server started on %s", port)
	log.Fatal(app.Listen(port))
}

// Handler untuk upload file
func uploadFile(c *fiber.Ctx) error {
	// Ambil file dari form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Cek ukuran file
	if file.Size > maxFileSize {
		return c.Status(400).JSON(fiber.Map{
			"error": "File size exceeds 100MB limit",
		})
	}

	// Generate nama file berdasarkan unix timestamp (sekarang + 1 jam) + ekstensi
	expiryTime := time.Now().Add(fileExpiry)

	// Ambil ekstensi dari file asli
	originalExt := filepath.Ext(file.Filename)
	if originalExt == "" {
		originalExt = ".bin" // default extension jika tidak ada
	}

	filename := strconv.FormatInt(expiryTime.Unix(), 10) + originalExt

	// Simpan file dengan nama unix timestamp + ekstensi
	filePath := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	return c.JSON(fiber.Map{
		"message":       "File uploaded successfully",
		"filename":      filename,
		"original_name": file.Filename,
		"size":          file.Size,
		"expires_at":    expiryTime.Format(time.RFC3339),
		"download_url":  fmt.Sprintf("/%s", filename),
	})
}

// Handler untuk download file
func downloadFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filePath := filepath.Join(uploadDir, filename)

	// Cek apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Parse unix timestamp dari nama file untuk validasi expiry
	// Hapus ekstensi untuk mendapatkan unix timestamp
	filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	unixTime, err := strconv.ParseInt(filenameWithoutExt, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid filename format",
		})
	}

	// Cek apakah file sudah expired
	expiryTime := time.Unix(unixTime, 0)
	if time.Now().After(expiryTime) {
		// Hapus file yang sudah expired
		os.Remove(filePath)
		return c.Status(404).JSON(fiber.Map{
			"error": "File has expired",
		})
	}

	// Download file
	return c.SendFile(filePath)
}

// Goroutine untuk membersihkan file yang sudah expired
func startCleanupRoutine() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredFiles()
		}
	}
}

// Fungsi untuk membersihkan file yang sudah expired
func cleanupExpiredFiles() {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Printf("Error reading upload directory: %v", err)
		return
	}

	now := time.Now()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Parse unix timestamp dari nama file
		filename := file.Name()
		filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
		unixTime, err := strconv.ParseInt(filenameWithoutExt, 10, 64)
		if err != nil {
			// Skip file yang bukan format unix timestamp
			continue
		}

		// Cek apakah file sudah expired
		expiryTime := time.Unix(unixTime, 0)
		if now.After(expiryTime) {
			filePath := filepath.Join(uploadDir, filename)
			if err := os.Remove(filePath); err != nil {
				log.Printf("Error removing expired file %s: %v", filename, err)
			} else {
				log.Printf("Removed expired file: %s", filename)
			}
		}
	}
}
