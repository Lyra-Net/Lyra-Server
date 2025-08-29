package router

import (
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
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	return r
}
