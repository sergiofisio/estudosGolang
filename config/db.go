package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"

)

func Connect() *sql.DB {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")

    connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    return db
}