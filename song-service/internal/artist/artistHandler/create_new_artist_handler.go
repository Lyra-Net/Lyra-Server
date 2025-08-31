package artistHandler

import (
	"log"
	"net/http"
	"song-service/dto"
	mq "song-service/internal/MQ"
	"song-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func CreateArtist(q *repository.Queries, producer *mq.KafkaProducer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var createArtistRequest dto.CreateArtistRequest
		if err := c.ShouldBind(&createArtistRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		artist, err := q.CreateArtist(c.Request.Context(), createArtistRequest.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = producer.Emit(c.Request.Context(), "artist_created", struct {
			ID   int32  `json:"id"`
			Name string `json:"name"`
		}{
			ID:   artist.ID,
			Name: artist.Name,
		})

		if err != nil {
			log.Println("failed to write message: %w", err)
		}
		log.Println("Event emited to Kafka")

		c.JSON(http.StatusCreated, artist.ID)
	}
}
