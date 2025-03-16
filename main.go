package main

import (
	"log"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/routes"
	"github.com/gofiber/fiber/v2"
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

	app := fiber.New()
	routes.Setup(app)

	app.Static("/", "./views")

	log.Fatal(app.Listen(":3000"))
}
