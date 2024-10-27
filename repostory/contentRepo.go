package repository

import (
	"rr/domain"

	"gorm.io/gorm"
)

type ContentRepo struct {
	DB *gorm.DB
}

func (r *ContentRepo) Create(content *domain.Content) error {
	return r.DB.Create(content).Error
}

func (r *ContentRepo) FindByID(id uint) (*domain.Content, error) {
	var content domain.Content
	if err := r.DB.First(&content, id).Error; err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *ContentRepo) Update(id uint, content *domain.Content) error {
	return r.DB.Save(content).Error
}

func (r *ContentRepo) Delete(id uint) error {
	return r.DB.Delete(&domain.Content{}, id).Error
}
