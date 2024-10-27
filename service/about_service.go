package service

import (
	"rr/domain"
	repository "rr/repostory"
)

type AboutService struct {
	Repo *repository.AboutRepo
}

func (s *AboutService) Create(t *domain.About) error {
	if t.Content == "" {
		t.Content = "Barada"
	}
	return s.Repo.Create(t)
}
func (s *AboutService) GetAll() ([]domain.About, error) {
	return s.Repo.FindAll()
}
func (s *AboutService) GetByID(id uint) (*domain.About, error) {
	return s.Repo.FindByID(id)
}
func (s *AboutService) Update(id uint, updatedAbout domain.About) (*domain.About, error) {
	existingAbout, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if updatedAbout.Content != "" {
		existingAbout.Content = updatedAbout.Content
	}

	if err := s.Repo.Update(id, existingAbout); err != nil {
		return nil, err
	}

	return existingAbout, nil
}
