package database

import (
	"bookstore-api/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	DB.AutoMigrate(&models.Book{}, &models.Category{}, &models.User{}, &models.Wishlist{})
}