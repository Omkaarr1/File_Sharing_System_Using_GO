package utils

import (
    "context"
    "log"
    "os"
    "github.com/jackc/pgx/v4/pgxpool"
)

func ConnectDB() *pgxpool.Pool {
    dbURL := os.Getenv("DATABASE_URL")
    conn, err := pgxpool.Connect(context.Background(), dbURL)
    if err != nil {
        log.Fatal("Unable to connect to database:", err)
    }
    return conn
}
