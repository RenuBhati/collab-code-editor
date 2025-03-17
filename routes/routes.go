package routes

import (
	"github.com/RenuBhati/editor/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/files", controllers.CreateFiles)
	app.Get("/files", controllers.ListFiles)
	app.Get("/files/:id", controllers.GetFile)
}
