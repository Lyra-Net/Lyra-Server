package artistHandler

import (
	"net/http"
	"song-service/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func ListArtist(q *repository.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		name := c.DefaultQuery("name", "")
		valid := true
		if name == "" {
			valid = false
		}
		artists, err := q.ListArtists(
			c.Request.Context(),
			repository.ListArtistsParams{
				Name:   pgtype.Text{String: name, Valid: valid},
				Limit:  int32(limit),
				Offset: (int32(page) - 1) * int32(limit),
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, artists)
	}
}
