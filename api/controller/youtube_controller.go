package controller

import (
	"github.com/Gohelraj/youtube-search-api/api/model"
	"github.com/Gohelraj/youtube-search-api/api/service"
	er "github.com/Gohelraj/youtube-search-api/error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type YoutubeController interface {
	GetVideos(c *gin.Context)
	SearchVideos(c *gin.Context)
}

type youtubeController struct {
	youtubeService service.YoutubeService
}

func NewYoutubeController(s service.YoutubeService) YoutubeController {
	return youtubeController{
		youtubeService: s,
	}
}

func (y youtubeController) GetVideos(c *gin.Context) {
	limitQueryParam := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitQueryParam)
	if err != nil {
		er.SendError(c, er.ErrInvalidValueInLimit)
		return
	}
	offsetQueryParam := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetQueryParam)
	if err != nil {
		er.SendError(c, er.ErrInvalidValueInOffset)
		return
	}
	videos, err := y.youtubeService.GetVideos(limit, offset)
	if err != nil {
		er.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, videos)
}

func (y youtubeController) SearchVideos(c *gin.Context) {
	var searchRequest model.SearchVideosRequest
	if err := c.BindJSON(&searchRequest); err != nil {
		er.SendError(c, err)
		return
	}
	if searchRequest.SearchString == "" {
		er.SendError(c, er.ErrSearchStringRequired)
		return
	}
	videos, err := y.youtubeService.SearchVideos(searchRequest.SearchString)
	if err != nil {
		er.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, videos)
}
