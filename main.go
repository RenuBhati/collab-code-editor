package main

import (
	"log"
	"os"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	database.InitDB()
	if _, err := os.Stat("./repos"); os.IsNotExist(err) {
		os.MkdirAll("./repos", os.ModePerm)
	}

	app := fiber.New()
	routes.Setup(app)

	app.Static("/", "./views")

	log.Fatal(app.Listen(":8080"))
}
