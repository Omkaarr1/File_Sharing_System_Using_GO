package utils

import (
    "context"
    "testing"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/stretchr/testify/assert"
)

// TestConnectDB tests the database connection
func TestConnectDB(t *testing.T) {
    // Mocking the database connection is not straightforward, so we will check for error handling
    dbURL := "postgres://invalid:invalid@localhost:5432/invalid_db"
    conn, err := pgxpool.Connect(context.Background(), dbURL)

    // Assert that the connection failed because of the invalid connection string
    assert.NotNil(t, err, "Expected error for invalid database connection")
    assert.Nil(t, conn, "Expected nil connection for invalid database")
}

// TestConnectDBSuccess tests a successful connection (adjust this for real DB)
func TestConnectDBSuccess(t *testing.T) {
    // Assuming a valid connection string is used, the test should pass when connected to a real database.
    // Example below assumes you have a local PostgreSQL instance running.
    dbURL := "postgres://username:password@localhost:5432/test_db"
    conn, err := pgxpool.Connect(context.Background(), dbURL)

    assert.Nil(t, err, "Expected no error for valid database connection")
    assert.NotNil(t, conn, "Expected valid database connection")

    if conn != nil {
        conn.Close()
    }
}
