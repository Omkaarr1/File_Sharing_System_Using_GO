package models

import (
	"file-sharing-system/utils"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser inserts a new user into the database
func CreateUser(user User) error {
	db := utils.ConnectDB() // Connect to the database
	defer db.Close()        // Ensure the connection is closed after the operation

	// MySQL query to insert a new user
	query := "INSERT INTO users (email, password_hash) VALUES (?, ?)"
	_, err := db.Exec(query, user.Email, user.Password)
	return err
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (User, error) {
	db := utils.ConnectDB() // Connect to the database
	defer db.Close()        // Ensure the connection is closed after the operation

	// MySQL query to get a user by email
	query := "SELECT id, email, password_hash FROM users WHERE email = ?"

	var user User
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}
