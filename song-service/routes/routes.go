package routes

import (
	"song-service/controllers"
	"song-service/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.Use(middleware.UserContextMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	songGroup := router.Group("/songs")
	{
		songGroup.POST("", controllers.CreateSongHandler(db))
		songGroup.GET("", controllers.GetSongsHandler(db))
	}
	artistGroup := router.Group("artists")
	{
		artistGroup.POST("", controllers.CreateArtistHandler(db))
		artistGroup.GET("", controllers.GetArtistHandler(db))
	}
	playlistGroup := router.Group("/playlists")
	{
		playlistGroup.POST("", controllers.CreatePlaylistHandler(db))
		playlistGroup.GET("", controllers.GetUserPlaylistsHandler(db))
	}
}
