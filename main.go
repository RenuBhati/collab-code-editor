package main

import (
	"log"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.InitDB()

	app := fiber.New()
	routes.Setup(app)

	app.Static("/", "./views")

	log.Fatal(app.Listen(":8080"))
}
