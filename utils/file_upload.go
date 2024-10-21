package utils

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func UploadFile(c *fiber.Ctx, formField string, uploadDir string, newFileName string) (string, error) {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	file, err := c.FormFile(formField)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	newFilePath := filepath.Join(uploadDir, newFileName+ext)

	if err := c.SaveFile(file, newFilePath); err != nil {
		return "", err
	}

	return newFilePath, nil
}
