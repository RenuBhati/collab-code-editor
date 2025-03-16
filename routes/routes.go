package routes

import (
	"github.com/RenuBhati/editor/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/files", controllers.CreateFiles)
}
