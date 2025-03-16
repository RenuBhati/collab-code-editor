package main

import (
	"log"

	"github.com/RenuBhati/editor/database"
)

func main() {
	adapter := database.SQLiteAdapter{}
	db, err := adapter.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	database.DB = db
	if err := database.MigrateDB(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
