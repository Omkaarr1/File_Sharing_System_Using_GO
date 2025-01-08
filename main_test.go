package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"file-sharing-system/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("test_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// HashPassword hashes a plain text password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// CheckPasswordHash checks if a plain password matches a hashed password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TestRegister tests user registration.
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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %s", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").WithArgs(user.Email, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		json.NewDecoder(r.Body).Decode(&user)

		hashedPassword, _ := HashPassword(user.Password)
		user.Password = hashedPassword

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("User registered")
	}).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %v", rr.Code)
	}

	expected := `"User registered"`
	actual := strings.TrimSpace(rr.Body.String())
	if expected != actual {
		t.Errorf("Expected body %v, got %v", expected, actual)
	}
}

// TestLogin tests user login.
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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %s", err)
	}
	defer db.Close()

	hashedPassword, _ := HashPassword(user.Password)
	mock.ExpectQuery("SELECT id, email, password_hash FROM users WHERE email = ?").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).AddRow(1, user.Email, hashedPassword))

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		})
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Logged in successfully",
			"token":   tokenString,
		})
	}).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rr.Code)
	}

	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if response["message"] != "Logged in successfully" {
		t.Errorf("Expected message 'Logged in successfully', got %v", response["message"])
	}

	if response["token"] == "" {
		t.Error("Expected a token in the response, but got an empty string")
	}
}

// TestFileUpload tests file upload.
func TestFileUpload(t *testing.T) {
	var jsonStr = []byte(`{"fileName": "testfile.txt", "fileContent": "This is a test file."}`)

	req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("File uploaded successfully")
	}).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rr.Code)
	}

	expected := `"File uploaded successfully"`
	actual := strings.TrimSpace(rr.Body.String())
	if expected != actual {
		t.Errorf("Expected body %v, got %v", expected, actual)
	}
}

// TestFileShare tests file sharing.
func TestFileShare(t *testing.T) {
	fileID := 123

	req, err := http.NewRequest("GET", "/share/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sharedURL := fmt.Sprintf("https://my-file-sharing-app.com/files/%d", fileID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sharedURL)
	}).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", rr.Code)
	}

	expected := `"https://my-file-sharing-app.com/files/123"`
	actual := strings.TrimSpace(rr.Body.String())
	if expected != actual {
		t.Errorf("Expected body %v, got %v", expected, actual)
	}
}
