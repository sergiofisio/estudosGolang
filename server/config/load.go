package config

import (
	"fmt"
	"os"
	"server/models"

	"github.com/joho/godotenv"

)

var AppConfig models.Config

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }

    AppConfig = models.Config{
        DBHost:     getEnv("DB_HOST", ""),
        DBUser:     getEnv("DB_USER", ""),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", ""),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        DBTimezone: getEnv("DB_TIMEZONE", "UTC"),
        JWTSecret:  getEnv("JWT_SECRET", ""),
    }
}

func getEnv(key, defaultValue string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return defaultValue
    }
    return value
}