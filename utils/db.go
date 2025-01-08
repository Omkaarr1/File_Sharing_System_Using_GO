package utils

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/joho/godotenv"         // Import for loading environment variables
)

// ConnectDB establishes a connection to the MySQL database
func ConnectDB() *sql.DB {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get the database URL from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in environment variables")
	}

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to ping the database: %v", err)
	}

	log.Println("Connected to the MySQL database successfully!")
	return db
}
