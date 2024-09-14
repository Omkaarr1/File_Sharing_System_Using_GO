package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "file-sharing-system/models"
    "github.com/DATA-DOG/go-sqlmock"
)

// TestRegister tests the user registration handler
func TestRegister(t *testing.T) {
    user := models.User{
        Email:    "test@example.com",
        Password: "password123",
    }
    userData, _ := json.Marshal(user)

    req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userData))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()

    // Mock database connection and insert operation
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock database: %s", err)
    }
    defer db.Close()

    mock.ExpectExec("INSERT INTO users").WithArgs(user.Email, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))

    handler := http.HandlerFunc(Register)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %v", rr.Code)
    }

    expected := `"User registered"`
    if rr.Body.String() != expected {
        t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
    }
}

// TestLogin tests the user login handler
func TestLogin(t *testing.T) {
    user := models.User{
        Email:    "test@example.com",
        Password: "password123",
    }
    userData, _ := json.Marshal(user)

    req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(userData))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()

    // Mock database connection and query
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock database: %s", err)
    }
    defer db.Close()

    hashedPassword, _ := HashPassword(user.Password)
    mock.ExpectQuery("SELECT id, email, password FROM users WHERE email = ?").
        WithArgs(user.Email).
        WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, user.Email, hashedPassword))

    handler := http.HandlerFunc(Login)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %v", rr.Code)
    }

    expected := `"Logged in successfully"`
    if rr.Body.String() != expected {
        t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
    }
}
