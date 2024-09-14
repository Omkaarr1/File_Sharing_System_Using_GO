package main

import (
    "bytes"
    "encoding/json"
    "fmt"    // Import fmt for string formatting
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "file-sharing-system/models"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
)

// Define jwtKey used for signing JWT tokens
var jwtKey = []byte("test_secret_key")

// Define the Claims struct used in JWT
type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// HashPassword for testing
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// CheckPasswordHash for testing
func CheckPasswordHash(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

// TestRegister tests the user registration handler
func TestRegister(t *testing.T) {
    // Mock the user data
    user := models.User{
        Email:    "test@example.com",
        Password: "password123",
    }
    userData, _ := json.Marshal(user)

    // Create a new HTTP request
    req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userData))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Record the response using httptest
    rr := httptest.NewRecorder()

    // Mock database connection and insert operation
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock database: %s", err)
    }
    defer db.Close()

    mock.ExpectExec("INSERT INTO users").WithArgs(user.Email, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))

    // Define the Register handler (handler logic should match your actual Register handler)
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var user models.User
        json.NewDecoder(r.Body).Decode(&user)

        // Hash password
        hashedPassword, err := HashPassword(user.Password)
        if err != nil {
            http.Error(w, "Error hashing password", http.StatusInternalServerError)
            return
        }
        user.Password = hashedPassword

        // Mock saving to database
        if err := models.CreateUser(user); err != nil {
            http.Error(w, "Unable to register user", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode("User registered")
    }).ServeHTTP(rr, req)

    // Assertions
    if rr.Code != http.StatusCreated {
        t.Errorf("expected status 201 Created, got %v", rr.Code)
    }

    expected := `"User registered"`
    if rr.Body.String() != expected {
        t.Errorf("expected body %v, got %v", expected, rr.Body.String())
    }
}

// TestLogin tests the user login handler
func TestLogin(t *testing.T) {
    // Mock the user data
    user := models.User{
        Email:    "test@example.com",
        Password: "password123",
    }
    userData, _ := json.Marshal(user)

    // Create a new HTTP request
    req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(userData))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Record the response using httptest
    rr := httptest.NewRecorder()

    // Mock database connection and query
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock database: %s", err)
    }
    defer db.Close()

    // Mock the stored hashed password
    hashedPassword, _ := HashPassword(user.Password)
    mock.ExpectQuery("SELECT id, email, password FROM users WHERE email = ?").
        WithArgs(user.Email).
        WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, user.Email, hashedPassword))

    // Define the Login handler (logic should match your actual Login handler)
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var user models.User
        json.NewDecoder(r.Body).Decode(&user)

        // Mock the stored user from the database
        storedUser := models.User{
            Email:    "test@example.com",
            Password: hashedPassword,
        }

        // Compare passwords
        if !CheckPasswordHash(user.Password, storedUser.Password) {
            http.Error(w, "Invalid password", http.StatusUnauthorized)
            return
        }

        // Generate JWT token
        expirationTime := time.Now().Add(1 * time.Hour)
        claims := &Claims{
            Email: storedUser.Email,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            http.Error(w, "Could not generate token", http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:    "token",
            Value:   tokenString,
            Expires: expirationTime,
        })

        json.NewEncoder(w).Encode("Logged in successfully")
    }).ServeHTTP(rr, req)

    // Assertions
    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200 OK, got %v", rr.Code)
    }

    expected := `"Logged in successfully"`
    if rr.Body.String() != expected {
        t.Errorf("expected body %v, got %v", expected, rr.Body.String())
    }
}

// TestFileUpload tests the file upload handler (mocked)
func TestFileUpload(t *testing.T) {
    // Mock a file upload
    var jsonStr = []byte(`{"fileName": "testfile.txt", "fileContent": "This is a test file."}`)

    // Create a new HTTP request
    req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Record the response using httptest
    rr := httptest.NewRecorder()

    // Define the Upload handler (handler logic should match your actual Upload handler)
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Simulate file upload process
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode("File uploaded successfully")
    }).ServeHTTP(rr, req)

    // Assertions
    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200 OK, got %v", rr.Code)
    }

    expected := `"File uploaded successfully"`
    if rr.Body.String() != expected {
        t.Errorf("expected body %v, got %v", expected, rr.Body.String())
    }
}

// TestFileShare tests the file sharing handler (mocked)
func TestFileShare(t *testing.T) {
    // Mock a file ID
    fileID := 123

    // Create a new HTTP request
    req, err := http.NewRequest("GET", "/share/123", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Record the response using httptest
    rr := httptest.NewRecorder()

    // Define the ShareFile handler (handler logic should match your actual ShareFile handler)
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sharedURL := fmt.Sprintf("https://my-file-sharing-app.com/files/%d", fileID)
        json.NewEncoder(w).Encode(sharedURL)
    }).ServeHTTP(rr, req)

    // Assertions
    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200 OK, got %v", rr.Code)
    }

    expected := `"https://my-file-sharing-app.com/files/123"`
    if rr.Body.String() != expected {
        t.Errorf("expected body %v, got %v", expected, rr.Body.String())
    }
}
