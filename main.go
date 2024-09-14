package main

import (
    "log"
    "net/http"

    "file-sharing-system/handlers"
    "file-sharing-system/utils"

    "github.com/gorilla/mux"
)

func main() {
    // Initialize database connection and Redis
    db := utils.ConnectDB()
    defer db.Close()
    
    redisClient := utils.ConnectRedis()
    defer redisClient.Close()

    // Initialize routes
    r := mux.NewRouter()

    // Auth routes
    r.HandleFunc("/register", handlers.Register).Methods("POST")
    r.HandleFunc("/login", handlers.Login).Methods("POST")

    // File routes
    r.HandleFunc("/upload", handlers.UploadFile).Methods("POST")
    r.HandleFunc("/files", handlers.GetFiles).Methods("GET")
    r.HandleFunc("/share/{file_id}", handlers.ShareFile).Methods("GET")
    
    // Start the server
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
