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

type NewsHandler struct {
	Service *service.NewsService
}

func (h *NewsHandler) Create(c *fiber.Ctx) error {
	var news domain.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}

	newFileName := "news_" + time.Now().Format("20060102150405")
	imagePath, err := utils.UploadFile(c, "news", "uploads/news", newFileName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Surat ýüklenip bilinmedi"})
	}
	news.Image = imagePath

	if err := h.Service.Create(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	news.IsActive = false
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	news.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, news.Image)
	return c.Status(fiber.StatusCreated).JSON(news)
}

func (h *NewsHandler) GetPaginated(c *fiber.Ctx) error {
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

	news, total, err := h.Service.GetPaginated(pageInt, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "newsler ýüklenip bilinmedi"})
	}

	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	for i := range news {
		news[i].Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, news[i].Image)
	}

	return c.JSON(fiber.Map{
		"data":  news,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

func (h *NewsHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	news, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "news tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "news tapylmady"})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	news.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, news.Image)
	return c.JSON(news)
}

func (h *NewsHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	news, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "news tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "news tapylmady"})
	}

	if news.Image != "" {
		t := news.Image
		fmt.Println("Pozuljak faýl:", t)
		err := utils.DeleteFileWithRetry(t)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Surat pozulyp bilinmedi: %v", err)})
		} else {
			fmt.Println("Surat üstünlikli pozuldy.")
		}
	}
	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "news pozulyp bilinmedi"})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "news üstünlikli pozuldy"})
}

func (h *NewsHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	news, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "news tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "news tapylmady"})
	}

	var updatedNews domain.News
	if err := c.BodyParser(&updatedNews); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	if file, err := c.FormFile("news"); err == nil && file != nil {
		if news.Image != "" {
			if err := os.Remove(news.Image); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Köne surat pozulyp bilinmedi"})
			}
		}
		newFileName := fmt.Sprintf("newsUpdate_%s", time.Now().Format("20060102150405"))
		imagePath, err := utils.UploadFile(c, "news", "uploads/news", newFileName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze surat ýüklenip bilinmedi"})
		}
		updatedNews.Image = imagePath

	} else {
		updatedNews.Image = news.Image
	}

	updatedNewsResult, err := h.Service.Update(uint(id), &updatedNews)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Habarlar tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Habarlar üýtgedilip bilinmedi"})
	}

	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	updatedNewsResult.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, updatedNews.Image)

	return c.JSON(updatedNewsResult)
}
