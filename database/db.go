package database

import (
	"fmt"
	"log"

	"github.com/RenuBhati/editor/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Database defines an interface
type DatabaseAdapter interface {
	Connect() (*gorm.DB, error)
}

type SQLiteAdapter struct{}

func (a SQLiteAdapter) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&models.File{})
	if err != nil {
		log.Println("Migration error:", err)
		return err
	}
	fmt.Println("Database migrated successfully")
	return nil
}
