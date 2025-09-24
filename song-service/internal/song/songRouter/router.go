package songrouter

import (
	mq "song-service/internal/MQ"
	"song-service/internal/repository"
	songhandler "song-service/internal/song/songHandler"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine, q *repository.Queries, producer *mq.KafkaProducer) {
	artistRoute := r.Group("/songs")
	{
		artistRoute.GET("", songhandler.ListSong(q))
		artistRoute.POST("", songhandler.CreateSong(q, producer))
		artistRoute.PUT("/:id", songhandler.UpdateSong(q, producer))
	}
}
