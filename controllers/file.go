package controllers

import (
	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/dto"
	"github.com/RenuBhati/editor/models"
	"github.com/gofiber/fiber/v2"
)

func CreateFiles(c *fiber.Ctx) error {
	var req dto.CreateFileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})

	}
	reqFile := models.File{
		Name:     req.Name,
		Content:  req.Content,
		OwnerID:  1,
		FileType: "owned",
	}
	myDb := database.DB
	if result := myDb.Create(&reqFile); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server side error"})

	}

}
