package handler

import (
	"rr/domain"
	"rr/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ContentHandler struct {
	Service *service.ContentService
}

func (h *ContentHandler) Create(c *fiber.Ctx) error {
	var content domain.Content
	if err := c.BodyParser(&content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}

	if err := h.Service.Create(&content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Maglumat saklanyp bilinmedi"})
	}

	return c.Status(fiber.StatusCreated).JSON(content)
}

func (h *ContentHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	content, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Maglumat tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Maglumat getirilip bilinmedi"})
	}

	return c.JSON(content)
}

func (h *ContentHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	var updatedContent domain.Content
	if err := c.BodyParser(&updatedContent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	content, err := h.Service.Update(uint(id), &updatedContent)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Maglumat tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Maglumat üýtgedilip bilinmedi"})
	}

	return c.JSON(content)
}

func (h *ContentHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Maglumat pozulyp bilinmedi"})
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
