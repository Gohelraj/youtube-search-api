package controller

import (
	"github.com/Gohelraj/youtube-search-api/api/service"
	"github.com/gin-gonic/gin"
)

type YoutubeController interface {
	SearchYoutube(c *gin.Context)
}

type youtubeController struct {
	youtubeService service.YoutubeService
}

func NewYoutubeController(s service.YoutubeService) YoutubeController {
	return youtubeController{
		youtubeService: s,
	}
}

func (y youtubeController) SearchYoutube(c *gin.Context) {

}
