package main

import (
	"os"
	"rr/database"
	"rr/routes"
	"rr/setup"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Database connection
	database.ConnectDB()

	// Fiber app setup with increased body limit
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024, // 500MB limit
	})
	app.Use(logger.New())

	routes.AuthRoutes(app)

	// Setup services for different resources
	HandlerBanner := setup.SetupServices(database.DB)
	HandlerEmployer := setup.SetupEmployerServices(database.DB)
	HandlerNews := setup.SetupNewsServices(database.DB)
	HandlerMedia := setup.SetupMediaServices(database.DB)
	HandlerLaws := setup.SetupLaws(database.DB)
	HandlerAbout := setup.SetupAboutServices(database.DB)
	HandlerContent := setup.SetupContenServices(database.DB)
	routes.SetupRoutes(
		app,
		HandlerBanner,
		HandlerEmployer,
		HandlerNews,
		HandlerMedia,
		HandlerLaws,
		HandlerAbout,
		HandlerContent,
	)
	ip := os.Getenv("BASE_URL")
	// Start server on port 5000
	port := os.Getenv("PORTW")
	app.Listen(ip + port)
}
