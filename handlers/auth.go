package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "file-sharing-system/models"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// CheckPasswordHash compares the hashed password with the plain one
func CheckPasswordHash(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

func Register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)
    
    // Hash password and store user in DB
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = hashedPassword

    // Save user to database
    if err := models.CreateUser(user); err != nil {
        http.Error(w, "Unable to register user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode("User registered")
}

func Login(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)
    
    storedUser, err := models.GetUserByEmail(user.Email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Compare the hashed passwords
    if !CheckPasswordHash(user.Password, storedUser.Password) {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    // Generate JWT Token
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
}
