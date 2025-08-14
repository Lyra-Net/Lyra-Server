package controllers

import (
	"net/http"
	"song-service/models"
	"song-service/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateSongInput struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	TitleToken []string `json:"title_token"`
	Categories []string `json:"categories"`
	ArtistIDs  []uint   `json:"artist_ids"`
}

func CreateSongHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateSongInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var artists []models.Artist
		if len(input.ArtistIDs) > 0 {
			if err := db.Find(&artists, input.ArtistIDs).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		song := models.Song{
			ID:         input.ID,
			Title:      input.Title,
			TitleToken: input.TitleToken,
			Categories: input.Categories,
			Artists:    artists,
		}

		if err := db.Create(&song).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, song)
	}
}

func GetSongsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		songs, total, err := services.GetSongs(db, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  songs,
			"page":  page,
			"limit": limit,
			"total": total,
		})
	}
}
