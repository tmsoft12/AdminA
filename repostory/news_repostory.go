package repository

import (
	"rr/domain"

	"gorm.io/gorm"
)

type NewsRepo struct {
	DB *gorm.DB
}

func (r *NewsRepo) Create(t *domain.News) error {
	return r.DB.Create(t).Error
}
func (r *NewsRepo) FindAll() ([]domain.News, error) {
	var news []domain.News
	err := r.DB.Find(&news).Error
	return news, err
}
func (r *NewsRepo) FindByID(id uint) (*domain.News, error) {
	var news domain.News
	err := r.DB.First(&news, id).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}
func (r *NewsRepo) Delete(id uint) error {
	return r.DB.Delete(&domain.News{}, id).Error
}
func (r *NewsRepo) Update(id uint, News *domain.News) error {
	return r.DB.Save(News).Error
}
