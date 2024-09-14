package models

import (
    "context"
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

func SaveFileMetadata(file File) error {
    db := utils.ConnectDB()
    defer db.Close()

    _, err := db.Exec(context.Background(), "INSERT INTO files (name, size, url, upload_date) VALUES ($1, $2, $3, $4)", file.Name, file.Size, file.URL, file.UploadDate)
    return err
}

func GetAllFiles() ([]File, error) {
    db := utils.ConnectDB()
    defer db.Close()

    rows, err := db.Query(context.Background(), "SELECT id, name, size, url, upload_date FROM files")
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
func GetFileByID(fileID string) (File, error) {
    db := utils.ConnectDB()
    defer db.Close()

    var file File
    err := db.QueryRow(context.Background(), "SELECT id, name, size, url, upload_date FROM files WHERE id = $1", fileID).Scan(&file.ID, &file.Name, &file.Size, &file.URL, &file.UploadDate)
    if err != nil {
        return File{}, err
    }
    return file, nil
}
