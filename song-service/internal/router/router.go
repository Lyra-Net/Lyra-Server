package router

import (
	"net/http"
	mq "song-service/internal/MQ"
	"song-service/internal/artist/artistRouter"
	"song-service/internal/repository"
	songrouter "song-service/internal/song/songRouter"

	"github.com/gin-gonic/gin"
)

func NewRouter(q *repository.Queries, producer *mq.KafkaProducer) *gin.Engine {
	r := gin.Default()

	artistRouter.Init(r, q, producer)
	songrouter.Init(r, q, producer)

	r.GET("/enums", func(c *gin.Context) {
		ctx := c.Request.Context()

		genres, err := q.GetEnumValues(ctx, "genre_enum")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		moods, err := q.GetEnumValues(ctx, "mood_enum")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"genres": genres,
			"moods":  moods,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}
