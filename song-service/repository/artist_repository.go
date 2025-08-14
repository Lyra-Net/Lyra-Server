package repository

import (
	"song-service/models"

	"gorm.io/gorm"
)

func CreateArtist(db *gorm.DB, artist *models.Artist) error {
	return db.Create(artist).Error
}

func GetArtists(db *gorm.DB, page, limit int, name string) ([]models.Artist, int64, error) {
	var artists []models.Artist
	var total int64

	query := db.Model(&models.Artist{})

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Songs").
		Limit(limit).
		Offset(offset).
		Find(&artists).Error

	return artists, total, err
}
