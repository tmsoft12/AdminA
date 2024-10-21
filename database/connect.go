package database

import (
	"log"
	"os"
	"rr/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	err = database.AutoMigrate(
		&domain.Banner{}, &domain.Employer{}, &domain.News{}, &domain.Media{}, &domain.User{}, &domain.Laws{},
	)
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}

	DB = database
	log.Println("Successfully connected to PostgreSQL")
}
