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

type EmployerHandler struct {
	Service *service.EmployerService
}

// Employer döretmek üçin funksiýa

func (h *EmployerHandler) Create(c *fiber.Ctx) error {
	var employer domain.Employer
	if err := c.BodyParser(&employer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}
	// Täze faýl adyny döretmek
	newFileName := "employer_" + time.Now().Format("20060102150405")
	imagePath, err := utils.UploadFile(c, "employer", "uploads/employer", newFileName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Surat ýüklenip bilinmedi"})
	}
	employer.Image = imagePath
	if err := h.Service.Create(&employer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	employer.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, employer.Image)
	return c.Status(fiber.StatusCreated).JSON(employer)

}

// Isgarleri sahypa boýunça getirmek
func (h *EmployerHandler) GetPaginated(c *fiber.Ctx) error {
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

	employer, total, err := h.Service.GetPaginated(pageInt, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Isgarler ýüklenip bilinmedi"})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	// Her bir isgarin surat URL-ni düzetmek
	for i := range employer {
		employer[i].Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, employer[i].Image)
	}

	return c.JSON(fiber.Map{
		"data":  employer,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

// ID boýunça isgar getirmek
func (h *EmployerHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nädogry ID"},
		)
	}

	employer, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Isgarler tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Isgarler tapylmady"})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	employer.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, employer.Image)
	return c.JSON(employer)
}

// Isgar pozmak üçin funksiýa
// Isgar pozmak üçin funksiýa
func (h *EmployerHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	employer, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Isgar tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Isgar tapylmady"})
	}

	// Isgar pozmak
	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Isgar pozulyp bilinmedi"})
	}

	// Suraty pozmak
	if employer.Image != "" {
		t := employer.Image
		fmt.Println("Pozuljak faýl:", t)
		err := utils.DeleteFileWithRetry(t)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Surat pozulyp bilinmedi: %v", err)})
		} else {
			fmt.Println("Surat üstünlikli pozuldy.")
		}
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Isgar üstünlikli pozuldy"})
}

// Isgar üýtgetmek üçin funksiýa
func (h *EmployerHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	employer, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Isgarler tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Isgarler tapylmady"})
	}

	var updatedEmployer domain.Employer
	if err := c.BodyParser(&updatedEmployer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	// Täze surat ýüklenipmi?
	if file, err := c.FormFile("employer"); err == nil && file != nil {
		// Köne suraty pozmak
		if employer.Image != "" {
			if err := os.Remove(employer.Image); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Köne surat pozulyp bilinmedi"})
			}
		}
		// Täze surat ýüklemek
		newFileName := fmt.Sprintf("employerUpdate_%s", time.Now().Format("20060102150405"))
		imagePath, err := utils.UploadFile(c, "employer", "uploads/employer", newFileName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze surat ýüklenip bilinmedi"})
		}
		updatedEmployer.Image = imagePath
	} else {
		// Täze surat ýüklenmedik bolsa, köne suraty sakla
		updatedEmployer.Image = employer.Image
	}

	updatedEmployerResult, err := h.Service.Update(uint(id), &updatedEmployer)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Isgarler tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Isgarler üýtgedilip bilinmedi"})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	updatedEmployerResult.Image = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, updatedEmployer.Image)
	fmt.Println(updatedEmployerResult)
	return c.JSON(updatedEmployerResult)
}
