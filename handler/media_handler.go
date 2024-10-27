package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"rr/domain"
	"rr/service"
	"rr/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MediaHandler struct {
	Service *service.MediaService
}

func (h *MediaHandler) Create(c *fiber.Ctx) error {
	var media domain.Media
	if err := c.BodyParser(&media); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}

	// Video dosyasını yükleyin
	newFileName := "video_" + time.Now().Format("20060102150405")
	videoPath, err := utils.UploadFile(c, "video", "uploads/media/video", newFileName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Video ýüklenip bilinmedi"})
	}
	videoFormat := filepath.Ext(videoPath) // Dosya formatını al
	media.Video = newFileName + videoFormat

	// Cover dosyasını yükleyin
	newCoverName := "cover_" + time.Now().Format("20060102150405")
	coverPath, err := utils.UploadFile(c, "cover", "uploads/media/cover", newCoverName) // Tanımlamayı düzelt
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cover ýüklenip bilinmedi"})
	}
	coverFormat := filepath.Ext(coverPath) // Dosya formatını al
	media.Cover = newCoverName + coverFormat

	// Medya kaydını veritabanına ekleyin
	if err := h.Service.Create(&media); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	media.Video = fmt.Sprintf("http://%s:%s/video/%s", ip, port, media.Video)
	media.Cover = fmt.Sprintf("http://%s:%s/api/admin/uploads/media/cover/%s", ip, port, media.Cover)

	return c.Status(fiber.StatusCreated).JSON(media)
}

func (h *MediaHandler) GetPaginated(c *fiber.Ctx) error {
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
	media, total, err := h.Service.GetPaginated(pageInt, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media ýüklenip bilinmedi"})
	}
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	// Her bir isgarin surat URL-ni düzetmek
	for i := range media {
		media[i].Video = fmt.Sprintf("http://%s:%s/video/%s", ip, port, media[i].Video)
	}
	for i := range media {
		media[i].Cover = fmt.Sprintf("http://%s:%s/api/admin/uploads/media/cover/%s", ip, port, media[i].Cover)
	}
	return c.JSON(fiber.Map{
		"data":  media,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}
func (h *MediaHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nädogry ID"},
		)
	}

	media, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "mediA tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "mediA tapylmady"})
	}

	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	media.Video = fmt.Sprintf("http://%s:%s/video/%s", ip, port, media.Video)
	media.Cover = fmt.Sprintf("http://%s:%s/api/admin/uploads/media/cover/%s", ip, port, media.Cover)
	return c.JSON(media)
}
func (h *MediaHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	media, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media tapylmady"})
	}

	// Veritabanından kaydı sil
	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media pozulyp bilinmedi"})
	}

	// Video ve Cover dosyalarının tam yolu
	videoPath := fmt.Sprintf("uploads/media/video/%s", media.Video)
	coverPath := fmt.Sprintf("uploads/media/cover/%s", media.Cover)

	// Video dosyasını silme işlemi
	if err := utils.DeleteFileWithRetry(videoPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Video pozulyp bilinmedi: %v", err)})
	}

	// Cover dosyasını silme işlemi
	if err := utils.DeleteFileWithRetry(coverPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Cover pozulyp bilinmedi: %v", err)})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Media üstünlikli pozuldy"})
}
func (h *MediaHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	media, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media tapylmady"})
	}

	var updateMedia domain.Media
	if err := c.BodyParser(&updateMedia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	// Video dosyasını güncelleme
	if fileVideo, err := c.FormFile("video"); err == nil && fileVideo != nil {
		if media.Video != "" {
			if err := os.Remove("uploads/media/video/" + media.Video); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Köne video pozulyp bilinmedi"})
			}
		}
		newVideoName := fmt.Sprintf("VideoUpdate_%s", time.Now().Format("20060102150405"))
		videoPath, err := utils.UploadFile(c, "video", "uploads/media/video", newVideoName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze video ýüklenip bilinmedi"})
		}
		updateMedia.Video = newVideoName + filepath.Ext(videoPath)
	} else {
		updateMedia.Video = media.Video
	}

	// Cover dosyasını güncelleme
	if fileCover, err := c.FormFile("cover"); err == nil && fileCover != nil {
		if media.Cover != "" {
			if err := os.Remove("uploads/media/cover/" + media.Cover); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Köne cover pozulyp bilinmedi"})
			}
		}
		newCoverName := fmt.Sprintf("CoverUpdate_%s", time.Now().Format("20060102150405"))
		coverPath, err := utils.UploadFile(c, "cover", "uploads/media/cover", newCoverName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze cover ýüklenip bilinmedi"})
		}
		updateMedia.Cover = newCoverName + filepath.Ext(coverPath)
	} else {
		updateMedia.Cover = media.Cover
	}

	// Güncellemeyi veritabanına kaydet
	updateMediaResult, err := h.Service.Update(uint(id), &updateMedia)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media üýtgedilip bilinmedi"})
	}

	// Tam URL'leri ayarla
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	updateMediaResult.Video = fmt.Sprintf("http://%s:%s/video/%s", ip, port, updateMedia.Video)
	updateMediaResult.Cover = fmt.Sprintf("http://%s:%s/api/admin/uploads/media/cover/%s", ip, port, updateMedia.Cover)

	return c.JSON(updateMediaResult)
}
