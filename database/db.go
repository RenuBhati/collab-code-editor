package database

import (
	"log"

	"github.com/RenuBhati/editor/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := DB.AutoMigrate(&models.File{}); err != nil {
		log.Println("Error during migration:", err)
		return err
	}

	return nil
}
