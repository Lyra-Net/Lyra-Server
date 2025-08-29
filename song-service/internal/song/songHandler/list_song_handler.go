package songhandler

import (
	"net/http"
	"song-service/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListSong(q *repository.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		songs, err := q.ListSongsWithArtists(
			c.Request.Context(),
			repository.ListSongsWithArtistsParams{
				Limit:  int32(limit),
				Offset: int32(page-1) * int32(limit),
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, songs)
	}
}
