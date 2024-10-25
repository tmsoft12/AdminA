package webdomain

import "rr/domain"

type HomePage struct {
	Banners   []domain.Banner   `json:"banners"`
	Employers []domain.Employer `json:"employers"`
	Laws      []domain.Laws     `json:"laws"`
	Media     []domain.Media    `json:"media"`
	News      []domain.News     `json:"news"`
}
