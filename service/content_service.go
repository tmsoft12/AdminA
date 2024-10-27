package service

import (
	"rr/domain"
	repository "rr/repostory"
)

type ContentService struct {
	Repo *repository.ContentRepo
}

func (s *ContentService) Create(content *domain.Content) error {
	return s.Repo.Create(content)
}

func (s *ContentService) GetByID(id uint) (*domain.Content, error) {
	return s.Repo.FindByID(id)
}

func (s *ContentService) Update(id uint, updatedContent *domain.Content) (*domain.Content, error) {
	existingContent, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Güncelleme işlemi
	existingContent.Top = updatedContent.Top
	existingContent.Bottom = updatedContent.Bottom

	if err := s.Repo.Update(id, existingContent); err != nil {
		return nil, err
	}

	return existingContent, nil
}

func (s *ContentService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
