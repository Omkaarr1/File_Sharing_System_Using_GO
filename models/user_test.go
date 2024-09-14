package models

import (
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
)

// TestCreateUser tests the CreateUser function
func TestCreateUser(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock database: %s", err)
    }
    defer db.Close()

    user := User{
        Email:    "test@example.com",
        Password: "password123",
    }

    mock.ExpectExec("INSERT INTO users").WithArgs(user.Email, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))

    if err := CreateUser(user); err != nil {
        t.Errorf("Error creating user: %s", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}

// TestGetUserByEmail tests the GetUserByEmail function
func TestGetUserByEmail(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock database: %s", err)
    }
    defer db.Close()

    user := User{
        Email:    "test@example.com",
        Password: "password123",
    }

    mock.ExpectQuery("SELECT id, email, password FROM users WHERE email = ?").
        WithArgs(user.Email).
        WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, user.Email, user.Password))

    returnedUser, err := GetUserByEmail(user.Email)
    if err != nil {
        t.Errorf("Error retrieving user: %s", err)
    }

    if returnedUser.Email != user.Email {
        t.Errorf("Expected email %v, got %v", user.Email, returnedUser.Email)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}
