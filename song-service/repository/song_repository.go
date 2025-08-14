package repository

import (
	"song-service/models"

	"gorm.io/gorm"
)

func GetSongs(db *gorm.DB, page, limit int) ([]models.Song, int64, error) {
	var songs []models.Song
	var total int64

	db.Model(&models.Song{}).Count(&total)

	offset := (page - 1) * limit
	err := db.Preload("Artists").Limit(limit).Offset(offset).Find(&songs).Error

	return songs, total, err
}

func GetSongById(db *gorm.DB, id string) (models.Song, error) {
	var song models.Song
	err := db.First(&song, "id = ?", id).Error
	return song, err
}

func SaveSong(db *gorm.DB, song *models.Song) error {
	return db.Create(song).Error
}
