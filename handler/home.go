package handler

import (
	repository "rr/repostory"

	"github.com/gofiber/fiber/v2"
)

// GetMediaByLanguage - saýlanan dildäki media maglumatlaryny bermek üçin handler
func GetMediaByLanguage(c *fiber.Ctx) error {
	language := c.Query("language") // Query bilen dil almak
	if language == "" {
		language = "tm" // Eger dil ýok bolsa, "en" kesgitlenýär
	}

	media, err := repository.GetMediaByLanguage(language)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	if media == nil || len(media) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No media found for the specified language."})
	}

	return c.JSON(media)
}
