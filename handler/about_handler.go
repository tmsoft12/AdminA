package handler

import (
	"rr/domain"
	"rr/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AboutHandler struct {
	Service *service.AboutService
}

func (h *AboutHandler) Create(c *fiber.Ctx) error {
	var about domain.About
	if err := c.BodyParser(&about); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}

	// Veritabanına kaydet
	if err := h.Service.Create(&about); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Maglumat saklanyp bilinmedi"})
	}

	return c.Status(fiber.StatusCreated).JSON(about)
}
func (h *AboutHandler) GetByID(c *fiber.Ctx) error {
	// idStr := c.Params("id")
	idStr := "1"
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	about, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Maglumat tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Maglumat getirilip bilinmedi"})
	}

	return c.JSON(about)
}

func (h *AboutHandler) Update(c *fiber.Ctx) error {
	// idStr := c.Params("id")
	idStr := "1"
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	var updatedAbout domain.About
	if err := c.BodyParser(&updatedAbout); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	about, err := h.Service.Update(uint(id), updatedAbout)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Maglumat tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Maglumat üýtgedilip bilinmedi"})
	}

	return c.JSON(about)
}
