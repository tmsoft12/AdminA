package handler

import (
	"rr/domain"
	"rr/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LawsHandler struct {
	Service *service.LawsService
}

// Create - Täze Laws ýazgyny döretmek üçin handler
func (h *LawsHandler) Create(c *fiber.Ctx) error {
	var laws domain.Laws

	// Maglumatlary almak
	if err := c.BodyParser(&laws); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar: Serwerde ýalňyşlyk ýüze çykdy"})
	}

	// Täze Laws ýazgyny döretmek
	if err := h.Service.Create(&laws); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Laws ýazgy döretmekde ýalňyşlyk ýüze çykdy"})
	}

	// Başarili goşma işlemi
	return c.Status(fiber.StatusCreated).JSON(laws)
}

func (h *LawsHandler) GetPaginated(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	laws, total, err := h.Service.GetPaginated(pageInt, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Laws ýüklenip bilinmedi"})
	}

	return c.JSON(fiber.Map{
		"data":  laws,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}
func (h *LawsHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	laws, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Laws tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Laws tapylmady"})
	}

	return c.JSON(laws)
}

// Delete - ID bilen Laws ýazgyny öçürmek üçin handler
func (h *LawsHandler) Delete(c *fiber.Ctx) error {
	// ID'ni almak we parse etmek
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	// ID bilen Laws ýazgyny tapmak
	existingLaws, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Laws tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Laws maglumatlary almakda ýalňyşlyk ýüze çykdy"})
	}

	// Laws ýazgyny öçürmek
	if err := h.Service.Delete(existingLaws.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Laws pozulyp bilinmedi"})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Laws üstünlikli pozuldy"})
}
func (h *LawsHandler) Update(c *fiber.Ctx) error {
	// ID'ni almak we parse etmek
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	// ID bilen Laws ýazgyny tapmak
	existingLaws, err := h.Service.GetByID(uint(id)) // İki değer alıyoruz
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Laws tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Laws maglumatlary almakda ýalňyşlyk ýüze çykdy"})
	}

	// Täze maglumatlary almak
	var updatedLaws domain.Laws
	if err := c.BodyParser(&updatedLaws); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	// Boş olmayan alanlarla mevcut kaydı güncelle
	if updatedLaws.Laws != "" {
		existingLaws.Laws = updatedLaws.Laws
	}

	if updatedLaws.Title != "" {
		existingLaws.Title = updatedLaws.Title
	}

	// Üýtgedilen Laws ýazgyny database'e ýazmak
	if err := h.Service.Update(uint(id), existingLaws); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Laws üýtgedilip bilinmedi"})
	}

	// Güncellenen kaydı döndür
	return c.JSON(existingLaws)
}
