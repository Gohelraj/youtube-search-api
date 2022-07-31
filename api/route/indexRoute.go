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
	// Default gin router with the Logger and Recovery middleware already attached.
	router := gin.Default()

	// Health check endpoint for the API server to check if it is running
	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello There!!")
		return
	})

	videoRepository := repository.NewVideoRepo(pgxPool)
	videoService := service.NewVideoService(videoRepository)
	videoController := controller.NewVideoController(videoService)

	// Videos API routes
	router.GET("/videos", videoController.GetVideos)
	router.POST("/videos/search", videoController.SearchVideos)

	return router
}
