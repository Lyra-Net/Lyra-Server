package songhandler

import (
	"net/http"
	"song-service/dto"
	mq "song-service/internal/MQ"
	"song-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func CreateSong(q *repository.Queries, producer *mq.KafkaProducer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dto.CreateSongRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		params := repository.CreateSongParams{
			ID:         req.ID,
			Title:      req.Title,
			TitleToken: req.TitleToken,
			Categories: req.Categories,
		}

		song, err := q.CreateSong(c.Request.Context(), params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, artistId := range req.ArtistIDS {
			err = q.AddSongArtists(c.Request.Context(), repository.AddSongArtistsParams{
				ArtistID: int32(artistId),
				SongID:   song.ID,
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		_ = producer.Emit(c.Request.Context(), "song_created", req)

		c.JSON(http.StatusCreated, gin.H{
			"id": song.ID,
		})
	}
}
