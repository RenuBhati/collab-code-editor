package main

import (
	"log"
	"os"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("error initializing database", err)
	}

	if _, err := os.Stat("./repos"); os.IsNotExist(err) {
		os.MkdirAll("./repos", os.ModePerm)
	}
	database.SeedDB()
	app := fiber.New()
	routes.Setup(app)

	app.Static("/", "./views")

	log.Fatal(app.Listen(":8080"))
}
