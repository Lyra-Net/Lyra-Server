package controllers

import (
	"net/http"
	"song-service/models"
	"song-service/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateArtistHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var artist models.Artist
		if err := c.ShouldBindJSON(&artist); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid Artist model - error": err.Error()})
			return
		}
		if err := services.CreateArtist(db, &artist); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, artist)
	}
}

func GetArtistHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		name := c.DefaultQuery("name", "")
		artist, total, err := services.GetArtists(db, page, limit, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":  artist,
			"page":  page,
			"limit": limit,
			"total": total,
		})
	}
}
