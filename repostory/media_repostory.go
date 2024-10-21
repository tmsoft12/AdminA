package repository

import (
	"rr/domain"
	"time"

	"gorm.io/gorm"
)

type MediaRepo struct {
	DB *gorm.DB
}

func (r *MediaRepo) Create(t *domain.Media) error {
	if t.Date == "" {
		t.Date = time.Now().Format("2006-01-02 15:04:05")
	}
	return r.DB.Create(t).Error
}
func (r *MediaRepo) FindAll() ([]domain.Media, error) {
	var media []domain.Media
	err := r.DB.Find(&media).Error
	return media, err
}

func (r *MediaRepo) FindByID(id uint) (*domain.Media, error) {
	var media domain.Media
	err := r.DB.First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}
func (r *MediaRepo) Delete(id uint) error {
	return r.DB.Delete(&domain.Media{}, id).Error
}
func (r *MediaRepo) Update(id uint, Media *domain.Media) error {
	return r.DB.Save(Media).Error
}
