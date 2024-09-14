package models

import (
    "context"
    "file-sharing-system/utils"
)

type User struct {
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func CreateUser(user User) error {
    db := utils.ConnectDB()
    defer db.Close()

    _, err := db.Exec(context.Background(), "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
    return err
}

func GetUserByEmail(email string) (User, error) {
    db := utils.ConnectDB()
    defer db.Close()

    var user User
    err := db.QueryRow(context.Background(), "SELECT id, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password)
    return user, err
}
