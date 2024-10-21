package repository

import (
	"rr/domain"

	"gorm.io/gorm"
)

type EmployerRepo struct {
	DB *gorm.DB
}

func (r *EmployerRepo) Create(t *domain.Employer) error {
	return r.DB.Create(t).Error
}
func (r *EmployerRepo) FindAll() ([]domain.Employer, error) {
	var employer []domain.Employer
	err := r.DB.Find(&employer).Error
	return employer, err
}
func (r *EmployerRepo) FindByID(id uint) (*domain.Employer, error) {
	var employer domain.Employer
	err := r.DB.First(&employer, id).Error
	if err != nil {
		return nil, err
	}
	return &employer, nil
}
func (r *EmployerRepo) Delete(id uint) error {
	return r.DB.Delete(&domain.Employer{}, id).Error
}
func (r *EmployerRepo) Update(id uint, Employer *domain.Employer) error {
	return r.DB.Save(Employer).Error
}
