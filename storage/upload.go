package storage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jerson2000/api-qfirst/models"
)

const (
	maxUploadSize = 100 << 20 // 100MB
	UploadDir     = "./uploads"
)

func isValidFile(filename string) bool {
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".pdf":  true,
		".txt":  true,
		".docs": true,
		".xlx":  true,
		".pptx": true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse multipart form
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		models.ResponseWithError(w, http.StatusBadRequest, "Error parsing multipart form")
		return
	}

	// Retrieve files from form
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		models.ResponseWithError(w, http.StatusBadRequest, "No files uploaded")
		return
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		models.ResponseWithError(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	var uploadedFiles []string

	// Process each uploaded file
	for _, fileHeader := range files {
		if fileHeader == nil {
			continue
		}

		// Validate file extension
		if !isValidFile(fileHeader.Filename) {
			log.Printf("Rejected file: %s", fileHeader.Filename)
			continue
		}

		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			continue
		}
		defer file.Close()

		// Save file to disk
		destPath := filepath.Join(UploadDir, filepath.Base(fileHeader.Filename))
		destFile, err := os.Create(destPath)
		if err != nil {
			log.Printf("Error saving file: %v", err)
			continue
		}
		defer destFile.Close()

		// Copy file content
		file.Seek(0, 0) // Reset file pointer
		if _, err := io.Copy(destFile, file); err != nil {
			log.Printf("Error writing file: %v", err)
			continue
		}

		// Add file URL to response
		uploadedFiles = append(uploadedFiles, fmt.Sprintf("http://localhost:3000/uploads/%s", filepath.Base(fileHeader.Filename)))
	}

	// Respond with uploaded file URLs
	if len(uploadedFiles) > 0 {
		payload := map[string]any{
			"message": "Files uploaded successfully",
			"files":   uploadedFiles,
		}

		models.ResponseWithJSON(w, http.StatusOK, payload)

	} else {
		models.ResponseWithError(w, http.StatusBadRequest, "No valid files uploaded")
	}
}

func ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(UploadDir)
	if err != nil {
		models.ResponseWithError(w, http.StatusInternalServerError, "Unable to read upload directory")
		return
	}

	var fileUrls []string
	for _, file := range files {
		if !file.IsDir() {
			fileUrls = append(fileUrls, fmt.Sprintf("http://localhost:3000/uploads/%s", file.Name()))
		}
	}

	if len(fileUrls) > 0 {
		payload := map[string]any{
			"files": fileUrls,
		}
		models.ResponseWithJSON(w, http.StatusOK, payload)
	} else {
		models.ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "No files found"})
	}
}
