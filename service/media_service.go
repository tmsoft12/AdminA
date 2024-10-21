package service

import (
	"rr/domain"
	repository "rr/repostory"
)

type MediaService struct {
	Repo *repository.MediaRepo
}

func (s *MediaService) Create(t *domain.Media) error {

	if t.TM_title == "" {
		t.TM_title = "Turkmence"
	}
	if t.EN_title == "" {
		t.EN_title = "Englis"
	}
	if t.RU_title == "" {
		t.RU_title = "Rusca"
	}

	return s.Repo.Create(t)
}
func (s *MediaService) GetAll() ([]domain.Media, error) {
	return s.Repo.FindAll()
}
func (s *MediaService) GetPaginated(page int, limit int) ([]domain.Media, int64, error) {
	var media []domain.Media
	var total int64

	if err := s.Repo.DB.Model(&domain.Media{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := s.Repo.DB.Offset(offset).Limit(limit).Find(&media).Error; err != nil {
		return nil, 0, err
	}

	return media, total, nil
}
func (s *MediaService) GetByID(id uint) (*domain.Media, error) {
	return s.Repo.FindByID(id)
}

func (s *MediaService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
func (s *MediaService) Update(id uint, updatedMedia *domain.Media) (*domain.Media, error) {
	// Köne media maglumatlaryny al
	existingMedia, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Boş däl bolsa, täzelik bilen çalyş
	if updatedMedia.TM_title != "" {
		existingMedia.TM_title = updatedMedia.TM_title
	}
	if updatedMedia.EN_title != "" {
		existingMedia.EN_title = updatedMedia.EN_title
	}
	if updatedMedia.RU_title != "" {
		existingMedia.RU_title = updatedMedia.RU_title
	}

	// Wideony täzeläň
	if updatedMedia.Video != "" {
		existingMedia.Video = updatedMedia.Video
	}

	// Cover suraty täzeläň
	if updatedMedia.Cover != "" {
		existingMedia.Cover = updatedMedia.Cover
	}

	// Täzelenen maglumatlary ýatda sakla
	if err := s.Repo.Update(id, existingMedia); err != nil {
		return nil, err
	}

	return existingMedia, nil
}
