package config

import (
	"log"
	"os"
	"server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

    dsn := "host=" + os.Getenv("DB_HOST") +
        " user=" + os.Getenv("DB_USER") +
        " password=" + os.Getenv("DB_PASSWORD") +
        " dbname=" + os.Getenv("DB_NAME") +
        " port=" + os.Getenv("DB_PORT") +
        " sslmode=" + os.Getenv("DB_SSLMODE") +
        " TimeZone=" + os.Getenv("DB_TIMEZONE")

    log.Println("Connecting to database with DSN:", dsn)

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database:", err)
    }

	err = database.AutoMigrate(&models.User{}, &models.Client{}, &models.Address{}, &models.Type{}, &models.Control{})
    if err != nil {
        log.Fatal("failed to migrate database:", err)
    }

    DB = database
    log.Println("Database connection established")

    controlTypes := []models.Type{
        {Name: "Sessão individual"},
        {Name: "Pacote"},
    }

    for _, controlType := range controlTypes {
        if err := DB.FirstOrCreate(&controlType, models.Type{Name: controlType.Name}).Error; err != nil {
            log.Fatal("failed to insert control type:", err)
        }
    }

    log.Println("Initial control types inserted")
}