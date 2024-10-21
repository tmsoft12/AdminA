package repository

import (
	"rr/domain"

	"gorm.io/gorm"
)

type BannerRepo struct {
	DB *gorm.DB
}

func (r *BannerRepo) Create(t *domain.Banner) error {
	return r.DB.Create(t).Error
}

func (r *BannerRepo) FindAll() ([]domain.Banner, error) {
	var banners []domain.Banner
	err := r.DB.Find(&banners).Error
	return banners, err
}
func (r *BannerRepo) FindByID(id uint) (*domain.Banner, error) {
	var banner domain.Banner
	err := r.DB.First(&banner, id).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}
func (r *BannerRepo) Delete(id uint) error {
	return r.DB.Delete(&domain.Banner{}, id).Error
}
func (r *BannerRepo) Update(id uint, banner *domain.Banner) error {
	return r.DB.Save(banner).Error
}
