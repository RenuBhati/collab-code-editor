package controllers

import (
	"github.com/RenuBhati/editor/dto"
	"github.com/gofiber/fiber/v2"
)

func CreateFiles(c *fiber.Ctx) error {
	var req dto.CreateFileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})

	}

}
