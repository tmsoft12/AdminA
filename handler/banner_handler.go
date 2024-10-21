package handler

import (
	"fmt"
	"os"
	"rr/domain"
	"rr/service"
	"rr/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BannerHandler struct {
	Service *service.BannerService
}

// Banner döretmek üçin funksiýa
func (h *BannerHandler) Create(c *fiber.Ctx) error {
	var banner domain.Banner
	if err := c.BodyParser(&banner); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}

	// Täze faýl adyny döretmek
	newFileName := "banner_" + time.Now().Format("20060102150405")
	imagePath, err := utils.UploadFile(c, "banner", "uploads/banner", newFileName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Surat ýüklenip bilinmedi"})
	}
	banner.Image = imagePath

	// Banner döretmek
	if err := h.Service.Create(&banner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	banner.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, banner.Image)
	return c.Status(fiber.StatusCreated).JSON(banner)
}

// Bannerleri sahypa boýunça getirmek
func (h *BannerHandler) GetPaginated(c *fiber.Ctx) error {
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

	banners, total, err := h.Service.GetPaginated(pageInt, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Bannerler ýüklenip bilinmedi"})
	}

	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	// Her bir banneriň surat URL-ni düzetmek
	for i := range banners {
		banners[i].Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, banners[i].Image)
	}

	return c.JSON(fiber.Map{
		"data":  banners,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

// ID boýunça banner getirmek
func (h *BannerHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	banner, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Banner tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Banner tapylmady"})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	banner.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, banner.Image)
	return c.JSON(banner)
}

// Banner pozmak üçin funksiýa
func (h *BannerHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	banner, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Banner tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Banner tapylmady"})
	}

	// Banneri pozmak
	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Banner pozulyp bilinmedi"})
	}

	// Suraty pozmak
	if banner.Image != "" {
		imagePath := banner.Image
		if err := os.Remove(imagePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Surat pozulyp bilinmedi"})
		}
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Banner üstünlikli pozuldy"})
}

// Banneri üýtgetmek üçin funksiýa
func (h *BannerHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	banner, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Banner tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Banner tapylmady"})
	}

	var updatedBanner domain.Banner
	if err := c.BodyParser(&updatedBanner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	// Täze surat ýüklenipmi?
	if file, err := c.FormFile("banner"); err == nil && file != nil {
		// Köne suraty pozmak
		if banner.Image != "" {
			if err := os.Remove(banner.Image); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Köne surat pozulyp bilinmedi"})
			}
		}
		// Täze surat ýüklemek
		newFileName := fmt.Sprintf("bannerUpdate_%s", time.Now().Format("20060102150405"))
		imagePath, err := utils.UploadFile(c, "banner", "uploads/banner", newFileName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze surat ýüklenip bilinmedi"})
		}
		updatedBanner.Image = imagePath
	} else {
		// Täze surat ýüklenmedik bolsa, köne suraty sakla
		updatedBanner.Image = banner.Image
	}

	updatedBannerResult, err := h.Service.Update(uint(id), &updatedBanner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Banner tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Banner üýtgedilip bilinmedi"})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	updatedBannerResult.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, updatedBanner.Image)
	return c.JSON(updatedBannerResult)
}
