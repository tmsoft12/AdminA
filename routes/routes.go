package routes

import (
	"rr/handler"
	"rr/middleware"
	"rr/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	BannerHandler *handler.BannerHandler,
	EmployerHandler *handler.EmployerHandler,
	NewsHandler *handler.NewsHandler,
	MediaHandler *handler.MediaHandler,
	LawsHandler *handler.LawsHandler,
	AboutHandler *handler.AboutHandler,
	ContentHandler *handler.ContentHandler,
) {
	app.Static("api/admin/uploads", "./uploads")
	app.Get("/video/:video", utils.Play)

	Admin := app.Group("api/admin/", middleware.JWTProtected())

	// Banners routes
	// Admin.Static("uploads", "./uploads")
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

	About := Admin.Group("about")
	About.Get("/:id", AboutHandler.GetByID)
	About.Post("/", AboutHandler.Create)
	About.Put("/:id", AboutHandler.Update)

	Content := Admin.Group("content")
	Content.Post("/", ContentHandler.Create)
	Content.Get("/:id", ContentHandler.GetByID)
	Content.Put("/:id", ContentHandler.Update)
}

func AuthRoutes(app *fiber.App) {
	app.Post("/register", handler.Register)
	app.Post("api/admin/login", handler.Login)
	app.Post("api/admin/logout", handler.Logout)

}

// func SetupHome(app *fiber.App) {
// 	Home := app.Group("/home")
// 	Home.Get("/", handler.GetMediaByLanguage)
// }
