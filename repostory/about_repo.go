package repository

import (
	"rr/domain"

	"gorm.io/gorm"
)

type AboutRepo struct {
	DB *gorm.DB
}

func (r *AboutRepo) Create(t *domain.About) error {
	return r.DB.Create(t).Error
}
func (r *AboutRepo) FindAll() ([]domain.About, error) {
	var about []domain.About
	err := r.DB.Find(&about).Error
	return about, err
}
func (r *AboutRepo) FindByID(id uint) (*domain.About, error) {
	var about domain.About
	err := r.DB.First(&about, id).Error
	if err != nil {
		return nil, err
	}
	return &about, nil
}
func (r *AboutRepo) Update(id uint, About *domain.About) error {
	return r.DB.Save(About).Error
}
