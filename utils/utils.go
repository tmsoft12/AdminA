package utils

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func Play(c *fiber.Ctx) error {
	fileName := c.Params("video")
	filePath := filepath.Join("./uploads/media/video/", fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found"})
	}
	return c.SendFile(filePath, false)
}
