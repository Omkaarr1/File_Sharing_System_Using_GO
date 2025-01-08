package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv" // For string-to-integer conversion
	"time"

	"github.com/gorilla/mux"
	"file-sharing-system/models" // This should correctly import your models package
	"file-sharing-system/utils"
)

// UploadFile handles file uploads
func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload to S3 or Local Storage
	fileURL, err := utils.UploadToS3(file, handler.Filename)
	if err != nil {
		http.Error(w, "Unable to upload file", http.StatusInternalServerError)
		return
	}

	// Save file metadata in the database
	fileMetadata := models.File{
		Name:       handler.Filename,
		Size:       handler.Size,
		URL:        fileURL,
		UploadDate: time.Now(),
	}
	if err := models.SaveFileMetadata(fileMetadata); err != nil {
		http.Error(w, "Error saving file metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully")
}

// GetFiles retrieves all files
func GetFiles(w http.ResponseWriter, r *http.Request) {
	// Retrieve all files for the authenticated user
	files, err := models.GetAllFiles()
	if err != nil {
		http.Error(w, "Unable to retrieve files", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(files)
}

// ShareFile generates a shareable URL for a file
func ShareFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileIDStr := vars["file_id"]

	// Convert fileID from string to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Retrieve file by ID
	file, err := models.GetFileByID(fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	sharedURL := fmt.Sprintf("https://my-file-sharing-app.com/files/%d", file.ID) // Use %d for integers
	json.NewEncoder(w).Encode(sharedURL)
}
