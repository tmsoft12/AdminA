package setup

import (
	"rr/handler"
	repository "rr/repostory"

	"rr/service"

	"gorm.io/gorm"
)

func SetupServices(db *gorm.DB) *handler.BannerHandler {
	bannerRepo := &repository.BannerRepo{DB: db}
	bannerService := &service.BannerService{Repo: bannerRepo}
	bannerHandler := &handler.BannerHandler{Service: bannerService}

	return bannerHandler
}

func SetupEmployerServices(db *gorm.DB) *handler.EmployerHandler {
	employerRepo := &repository.EmployerRepo{DB: db}
	employerService := &service.EmployerService{Repo: employerRepo}
	employerHandler := &handler.EmployerHandler{Service: employerService}

	return employerHandler
}

func SetupNewsServices(db *gorm.DB) *handler.NewsHandler {
	newsRepo := &repository.NewsRepo{DB: db}
	newsService := &service.NewsService{Repo: newsRepo}
	newsHandler := &handler.NewsHandler{Service: newsService}

	return newsHandler
}
func SetupMediaServices(db *gorm.DB) *handler.MediaHandler {
	mediaRepo := &repository.MediaRepo{DB: db}
	mediaService := &service.MediaService{Repo: mediaRepo}
	mediaHandler := &handler.MediaHandler{Service: mediaService}

	return mediaHandler
}
func SetupLaws(db *gorm.DB) *handler.LawsHandler {
	lawsRepo := &repository.LawsRepo{DB: db}
	lawsService := &service.LawsService{Repo: lawsRepo}
	lawsHandler := &handler.LawsHandler{Service: lawsService}
	return lawsHandler
}
func SetupAboutServices(db *gorm.DB) *handler.AboutHandler {
	aboutRepo := &repository.AboutRepo{DB: db}
	aboutService := &service.AboutService{Repo: aboutRepo}
	aboutHandler := &handler.AboutHandler{Service: aboutService}

	return aboutHandler
}
