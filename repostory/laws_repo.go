package repository

import (
	"rr/domain"

	"gorm.io/gorm"
)

type LawsRepo struct {
	DB *gorm.DB
}

func (l *LawsRepo) Create(d *domain.Laws) error {
	return l.DB.Create(d).Error
}
func (l *LawsRepo) FindAll() ([]domain.Laws, error) {
	var laws []domain.Laws
	err := l.DB.Find(&laws).Error
	return laws, err
}
func (l *LawsRepo) FindByID(id uint) (*domain.Laws, error) {
	var law domain.Laws
	err := l.DB.First(&law, id).Error
	if err != nil {
		return nil, err
	}
	return &law, err
}
func (l *LawsRepo) Delete(id uint) error {
	return l.DB.Delete(&domain.Laws{}, id).Error
}
func (l *LawsRepo) Update(id uint, Media *domain.Media) error {
	return l.DB.Save(Media).Error
}
