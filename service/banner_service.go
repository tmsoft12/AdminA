package service

import (
	"rr/domain"
	repository "rr/repostory"
)

type BannerService struct {
	Repo *repository.BannerRepo
}

func (s *BannerService) Create(t *domain.Banner) error {

	if t.Link == "" {
		t.Link = "http://test.com"
	}
	return s.Repo.Create(t)
}
func (s *BannerService) GetAll() ([]domain.Banner, error) {
	return s.Repo.FindAll()
}
func (s *BannerService) GetPaginated(page int, limit int) ([]domain.Banner, int64, error) {
	var banners []domain.Banner
	var total int64

	if err := s.Repo.DB.Model(&domain.Banner{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := s.Repo.DB.Offset(offset).Limit(limit).Find(&banners).Error; err != nil {
		return nil, 0, err
	}

	return banners, total, nil
}
func (s *BannerService) GetByID(id uint) (*domain.Banner, error) {
	return s.Repo.FindByID(id)
}

func (s *BannerService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
func (s *BannerService) Update(id uint, updatedBanner *domain.Banner) (*domain.Banner, error) {
	existingBanner, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if updatedBanner.Link != "" {
		existingBanner.Link = updatedBanner.Link
	}
	if updatedBanner.Image != "" {
		existingBanner.Image = updatedBanner.Image
	}
	if updatedBanner.Is_Active != nil {
		existingBanner.Is_Active = updatedBanner.Is_Active
	}

	if err := s.Repo.Update(id, existingBanner); err != nil {
		return nil, err
	}

	return existingBanner, nil
}
