package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
)

const (
	maxFileSize = 500 * 1024 * 1024 // 500MB max file size
)

var allowedVideoTypes = map[string]bool{
	"video/mp4":       true,
	"video/webm":      true,
	"video/avi":       true,
	"video/mov":       true,
	"video/quicktime": true,
}

var allowedPresentationTypes = map[string]bool{
	"application/pdf": true, // .pdf
}

// UploadCourseFile handles file uploads for course modules
func (as *Server) UploadCourseFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileType := vars["type"] // "video" or "presentation"
	
	if fileType != "video" && fileType != "presentation" {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid file type"}, http.StatusBadRequest)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "File too large"}, http.StatusRequestEntityTooLarge)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "No file provided"}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	contentType := fileHeader.Header.Get("Content-Type")
	if fileType == "video" && !allowedVideoTypes[contentType] {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid video format. Allowed: MP4, WebM, AVI, MOV"}, http.StatusBadRequest)
		return
	}
	if fileType == "presentation" && !allowedPresentationTypes[contentType] {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid presentation format. Only PDF files are allowed"}, http.StatusBadRequest)
		return
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s%s", timestamp, strings.ReplaceAll(fileHeader.Filename, ext, ""), ext)
	
	// Create upload path
	uploadDir := filepath.Join("uploads", "courses", fileType+"s")
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Failed to create upload directory"}, http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(uploadDir, filename)
	
	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Failed to create file"}, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file contents
	fileSize, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		JSONResponse(w, models.Response{Success: false, Message: "Failed to save file"}, http.StatusInternalServerError)
		return
	}

	// Return file info
	response := map[string]interface{}{
		"success":           true,
		"message":           "File uploaded successfully",
		"filename":          filename,
		"original_filename": fileHeader.Filename,
		"file_path":         filePath,
		"file_size":         fileSize,
		"mime_type":         contentType,
		"uploaded_date":     time.Now(),
	}

	JSONResponse(w, response, http.StatusOK)
}

// ServeCourseFile serves uploaded course files
func (as *Server) ServeCourseFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileType := vars["type"]   // "video" or "presentation"
	filename := vars["filename"]

	if fileType != "video" && fileType != "presentation" {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Construct file path
	filePath := filepath.Join("uploads", "courses", fileType+"s", filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set appropriate headers
	if fileType == "video" {
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
	} else if fileType == "presentation" {
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
		// Allow PDF files to be embedded in iframes from the same origin
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'self';")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// DeleteCourseFile deletes an uploaded file
func (as *Server) DeleteCourseFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileType := vars["type"]
	filename := vars["filename"]

	if fileType != "video" && fileType != "presentation" {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid file type"}, http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", "courses", fileType+"s", filename)
	
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			JSONResponse(w, models.Response{Success: false, Message: "File not found"}, http.StatusNotFound)
		} else {
			JSONResponse(w, models.Response{Success: false, Message: "Failed to delete file"}, http.StatusInternalServerError)
		}
		return
	}

	JSONResponse(w, models.Response{Success: true, Message: "File deleted successfully"}, http.StatusOK)
}

