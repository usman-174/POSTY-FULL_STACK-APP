package database

import (
	"fmt"
	"log"
	"os"

	"github.com/usman-174/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func ConnectDataBase() *gorm.DB {
	fmt.Println("database.go START")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.Like{})
	fmt.Println("database.go STOP")

	return db
}
