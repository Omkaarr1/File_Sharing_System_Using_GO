package models

import (
	"database/sql"
	"file-sharing-system/utils"
	"time"
)

type File struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	URL        string    `json:"url"`
	UploadDate time.Time `json:"upload_date"`
}

// SaveFileMetadata inserts a file's metadata into the database
func SaveFileMetadata(file File) error {
	db := utils.ConnectDB() // Connect to the database
	defer db.Close()        // Close the connection after the operation

	// MySQL query to insert file metadata
	query := "INSERT INTO files (name, size, url, upload_date) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, file.Name, file.Size, file.URL, file.UploadDate)
	return err
}

// GetAllFiles retrieves all files from the database
func GetAllFiles() ([]File, error) {
	db := utils.ConnectDB() // Connect to the database
	defer db.Close()        // Close the connection after the operation

	// MySQL query to fetch all files
	query := "SELECT id, name, size, url, upload_date FROM files"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var file File
		if err := rows.Scan(&file.ID, &file.Name, &file.Size, &file.URL, &file.UploadDate); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// GetFileByID retrieves a file by its ID from the database
func GetFileByID(fileID int) (File, error) {
	db := utils.ConnectDB() // Connect to the database
	defer db.Close()        // Close the connection after the operation

	// MySQL query to fetch file by ID
	query := "SELECT id, name, size, url, upload_date FROM files WHERE id = ?"
	var file File
	err := db.QueryRow(query, fileID).Scan(&file.ID, &file.Name, &file.Size, &file.URL, &file.UploadDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return File{}, nil // No file found
		}
		return File{}, err
	}
	return file, nil
}
