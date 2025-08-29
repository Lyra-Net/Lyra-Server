package artistRouter

import (
	mq "song-service/internal/MQ"
	"song-service/internal/artist/artistHandler"
	"song-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine, q *repository.Queries, producer *mq.KafkaProducer) {
	artistRoute := r.Group("/artists")
	{
		artistRoute.GET("", artistHandler.ListArtist(q))
		artistRoute.POST("", artistHandler.CreateArtist(q, producer))
	}
}
