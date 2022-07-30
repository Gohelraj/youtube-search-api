package route

import (
	"github.com/Gohelraj/youtube-search-api/api/controller"
	"github.com/Gohelraj/youtube-search-api/api/repository"
	"github.com/Gohelraj/youtube-search-api/api/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

// InitializeRouter initialize all API routes
func InitializeRouter(pgxPool *pgxpool.Pool) *gin.Engine {
	// default gin router with the Logger and Recovery middleware already attached.
	r := gin.Default()

	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello There!!")
		return
	})

	youtubeRepository := repository.NewYoutubeRepo(pgxPool)
	youtubeService := service.NewYoutubeService(youtubeRepository)
	youtubeController := controller.NewYoutubeController(youtubeService)

	r.GET("/videos", youtubeController.GetVideos)
	r.POST("/videos/search", youtubeController.SearchVideos)

	return r
}
