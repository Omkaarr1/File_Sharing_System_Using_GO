package utils

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// TestConnectDB tests the database connection with an invalid connection string
func TestConnectDB(t *testing.T) {
	// Invalid MySQL connection string
	dbURL := "invalid_user:invalid_pass@tcp(localhost:3306)/invalid_db"
	db, err := sql.Open("mysql", dbURL)

	// Ensure that sql.Open does not return an error for invalid strings
	assert.Nil(t, err, "Expected no error from sql.Open for invalid connection string")
	assert.NotNil(t, db, "Expected a non-nil DB object from sql.Open")

	// Test Ping() to validate the connection
	if db != nil {
		err = db.Ping()
		assert.NotNil(t, err, "Expected error when pinging with an invalid connection")
		db.Close()
	}
}


// TestConnectDBSuccess tests a successful connection to a MySQL database
func TestConnectDBSuccess(t *testing.T) {
	// Valid MySQL connection string
	dbURL := "root:0000@tcp(localhost:3306)/file_sharing_db"
	db, err := sql.Open("mysql", dbURL)

	// Assert that the connection is successful
	assert.Nil(t, err, "Expected no error for valid database connection")
	assert.NotNil(t, db, "Expected valid database connection")

	// Test if the connection is alive
	if db != nil {
		err = db.Ping()
		assert.Nil(t, err, "Expected no error when pinging the database")

		// Close the connection
		err = db.Close()
		assert.Nil(t, err, "Expected no error when closing the database connection")
	}
}
