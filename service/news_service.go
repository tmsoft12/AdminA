package service

import (
	"rr/domain"
	repository "rr/repostory"
	"time"
)

type NewsService struct {
	Repo *repository.NewsRepo
}

func (s *NewsService) Create(t *domain.News) error {

	if t.TM_title == "" {
		t.TM_title = "Türkmen dilinde başlyk"
	}
	if t.TM_description == "" {
		t.TM_description = "Türkmen dilinde düşündiriş"
	}
	if t.EN_title == "" {
		t.EN_title = "Iňlis dilinde başlyk"
	}
	if t.EN_description == "" {
		t.EN_description = "Iňlis dilinde düşündiriş"
	}
	if t.RU_title == "" {
		t.RU_title = "Rus dilinde başlyk"
	}
	if t.RU_description == "" {
		t.RU_description = "Rus dilinde düşündiriş"
	}
	if t.Date == "" {
		t.Date = time.Now().Format("2006-01-02 15:04:05")
	}

	return s.Repo.Create(t)
}

func (s *NewsService) GetAll() ([]domain.News, error) {
	return s.Repo.FindAll()
}

func (s *NewsService) GetPaginated(page int, limit int) ([]domain.News, int64, error) {
	var news []domain.News
	var total int64

	if err := s.Repo.DB.Model(&domain.News{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := s.Repo.DB.Offset(offset).Limit(limit).Find(&news).Error; err != nil {
		return nil, 0, err
	}

	return news, total, nil
}

func (s *NewsService) GetByID(id uint) (*domain.News, error) {
	return s.Repo.FindByID(id)
}

func (s *NewsService) Delete(id uint) error {
	return s.Repo.Delete(id)
}

func (s *NewsService) Update(id uint, updatedNews *domain.News) (*domain.News, error) {
	existingNews, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if updatedNews.Image != "" {
		existingNews.Image = updatedNews.Image
	}
	// Galan maglumatlary täzele
	if updatedNews.TM_title != "" {
		existingNews.TM_title = updatedNews.TM_title
	}
	if updatedNews.TM_description != "" {
		existingNews.TM_description = updatedNews.TM_description
	}
	if updatedNews.EN_title != "" {
		existingNews.EN_title = updatedNews.EN_title
	}
	if updatedNews.EN_description != "" {
		existingNews.EN_description = updatedNews.EN_description
	}
	if updatedNews.RU_title != "" {
		existingNews.RU_title = updatedNews.RU_title
	}
	if updatedNews.RU_description != "" {
		existingNews.RU_description = updatedNews.RU_description
	}

	// Maglumatlar bazasynda täzelenýär
	if err := s.Repo.Update(id, existingNews); err != nil {
		return nil, err
	}

	return existingNews, nil
}
