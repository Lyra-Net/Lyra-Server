package songhandler

import (
	"net/http"
	"song-service/dto"
	mq "song-service/internal/MQ"
	"song-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func UpdateSong(q *repository.Queries, producer *mq.KafkaProducer) func(c *gin.Context) {
	return func(c *gin.Context) {
		songID := c.Param("id")
		if songID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing song id"})
			return
		}

		var req dto.UpdateSongRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		params := repository.UpdateSongParams{
			ID:         songID,
			Title:      req.Title,
			TitleToken: req.TitleToken,
			Categories: req.Categories,
			Duration:   pgtype.Int4{Valid: req.Duration > 0, Int32: req.Duration},
			Genre:      repository.NullGenreEnum{Valid: req.Genre != "", GenreEnum: repository.GenreEnum(req.Genre)},
			Mood:       repository.NullMoodEnum{Valid: req.Mood != "", MoodEnum: repository.MoodEnum(req.Mood)},
		}

		song, err := q.UpdateSong(c.Request.Context(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Clear old artists
		err = q.RemoveSongArtists(c.Request.Context(), songID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Re-attach artists
		for _, artistId := range req.Artists {
			err = q.AddSongArtists(c.Request.Context(), repository.AddSongArtistsParams{
				ArtistID: artistId.ID,
				SongID:   song.ID,
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		req.ID = songID
		_ = producer.Emit(c.Request.Context(), "song_created", req)

		c.JSON(http.StatusOK, gin.H{
			"id": song.ID,
		})
	}
}
