package repository

import (
	"rr/database"
	"rr/domain"
)

// GetMediaByLanguage - belli bir dilde media maglumatlaryny almak üçin funksiýa
func GetMediaByLanguage(language string) ([]domain.Media, error) {
	var mediaList []domain.Media

	// Database'den maglumatlary almak
	err := database.DB.Where("tm_title IS NOT NULL OR en_title IS NOT NULL OR ru_title IS NOT NULL").Find(&mediaList).Error
	if err != nil {
		return nil, err
	}

	var filteredMedia []domain.Media
	for _, media := range mediaList {
		var selectedTitle string

		// Saýlanan dile görä başlygy saýlamak
		switch language {
		case "tm":
			selectedTitle = media.TM_title
		case "en":
			selectedTitle = media.EN_title
		case "ru":
			selectedTitle = media.RU_title
		default:
			selectedTitle = media.TM_title // Default dil hökmünde "tm"
		}

		// Diňe saýlanan başlygy görkezmek
		filteredMedia = append(filteredMedia, domain.Media{
			ID:       media.ID,
			Cover:    media.Cover,
			Video:    media.Video,
			Date:     media.Date,
			View:     media.View,
			TM_title: "", // Başga diller üçin başlyk boşaldyldy
			EN_title: "", // Başga diller üçin başlyk boşaldyldy
			RU_title: "", // Başga diller üçin başlyk boşaldyldy
		})

		// Seçilen başlygy berjaý etmek
		switch language {
		case "tm":
			filteredMedia[len(filteredMedia)-1].TM_title = selectedTitle
		case "en":
			filteredMedia[len(filteredMedia)-1].EN_title = selectedTitle
		case "ru":
			filteredMedia[len(filteredMedia)-1].RU_title = selectedTitle
		}
	}

	return filteredMedia, nil
}
