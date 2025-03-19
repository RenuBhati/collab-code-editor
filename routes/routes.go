package routes

import (
	"github.com/RenuBhati/editor/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/files", controllers.CreateFiles)
	app.Get("/files", controllers.ListFiles)
	app.Get("/files/:id", controllers.GetFile)
	app.Put("/files/:id", controllers.UpdateFiles)
	app.Delete("/files/:id", controllers.DeleteFiles)
	app.Post("/files/:id/share", controllers.ShareFile)
	app.Post("/files/:id/save", controllers.SaveFile)
	app.Get("/files/:id/history", controllers.FileHistory)
	app.Get("/files/:id/blame", controllers.GitBlame)
}
