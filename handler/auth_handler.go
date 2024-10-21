package handler

import (
	"rr/database"
	"rr/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	type Talap struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Talap
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry talap"})
	}

	err := service.RegisterUser(database.DB, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Ulanyjy üstünlikli hasaba alndy"})
}

func Login(c *fiber.Ctx) error {
	type Talap struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Talap
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry talap"})
	}

	token, err := service.LoginUser(database.DB, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Nädogry ulanyjy ýa-da parol"})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Girip bildiňiz"})
}

func Logout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour) // Kukini möhletini geçirmek
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Çykdyňyz"})
}
