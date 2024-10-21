package routes

import (
	"rr/handler"
	"rr/middleware"
	"rr/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, BannerHandler *handler.BannerHandler, EmployerHandler *handler.EmployerHandler, NewsHandler *handler.NewsHandler, MediaHandler *handler.MediaHandler, LawsHandler *handler.LawsHandler) {

	// General admin group with JWT protection
	Admin := app.Group("api/admin/", middleware.JWTProtected())

	// Banners routes
	Admin.Static("uploads", "./uploads")
	Admin.Post("banners", BannerHandler.Create)
	Admin.Get("banners", BannerHandler.GetPaginated)
	Admin.Get("banners/:id", BannerHandler.GetByID)
	Admin.Delete("banners/:id", BannerHandler.Delete)
	Admin.Put("banners/:id", BannerHandler.Update)

	// Employers routes
	Admin.Post("employers", EmployerHandler.Create)
	Admin.Get("employers/:id", EmployerHandler.GetByID)
	Admin.Get("employers", EmployerHandler.GetPaginated)
	Admin.Delete("employers/:id", EmployerHandler.Delete)
	Admin.Put("employers/:id", EmployerHandler.Update)

	// News routes
	Admin.Post("news", NewsHandler.Create)
	Admin.Get("news/:id", NewsHandler.GetByID)
	Admin.Get("news", NewsHandler.GetPaginated)
	Admin.Delete("news/:id", NewsHandler.Delete)
	Admin.Put("news/:id", NewsHandler.Update)

	// Media routes
	Media := Admin.Group("media")
	Media.Get("/video/:video", utils.Play)
	Media.Post("/", MediaHandler.Create)
	Media.Get("/:id", MediaHandler.GetByID)
	Media.Get("/", MediaHandler.GetPaginated)
	Media.Delete("/:id", MediaHandler.Delete)
	Media.Put("/:id", MediaHandler.Update)
	//Laws routes
	Laws := Admin.Group("laws")
	Laws.Post("/", LawsHandler.Create)
	Laws.Get("/:id", LawsHandler.GetByID)
	Laws.Get("/", LawsHandler.GetPaginated)
	Laws.Delete("/:id", LawsHandler.Delete)
	Laws.Put("/:id", LawsHandler.Update)
	// Protecting uploads with JWT middleware
	app.Static("/uploads", "./uploads", fiber.Static{
		Browse: true, // Optional: to allow browsing files in folder
		Next: func(c *fiber.Ctx) bool {
			// Call JWTProtected and check for error
			if err := middleware.JWTProtected()(c); err != nil {
				return true // Block access to static files
			}
			return false // Allow access to static files
		},
	})

}

func AuthRoutes(app *fiber.App) {
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)
	app.Post("/logout", handler.Logout)

	// Protected route, requires JWT authentication
	app.Get("/protected", middleware.JWTProtected(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "You are authorized"})
	})
}
func SetupHome(app *fiber.App) {
	Home := app.Group("/home")
	Home.Get("/", handler.GetMediaByLanguage)
}
