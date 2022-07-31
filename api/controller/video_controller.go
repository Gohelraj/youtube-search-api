package controller

import (
	"github.com/Gohelraj/youtube-search-api/api/model"
	"github.com/Gohelraj/youtube-search-api/api/service"
	er "github.com/Gohelraj/youtube-search-api/error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type VideoController interface {
	GetVideos(c *gin.Context)
	SearchVideos(c *gin.Context)
}

type videoController struct {
	videoService service.VideoService
}

func NewVideoController(s service.VideoService) VideoController {
	return videoController{
		videoService: s,
	}
}

// GetVideos returns videos from the database
func (v videoController) GetVideos(c *gin.Context) {
	limitQueryParam := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitQueryParam)
	if err != nil {
		er.SendError(c, er.ErrInvalidValueInLimit)
		return
	}
	if limit > 100 {
		er.SendError(c, er.ErrLimitExceeded)
		return
	}
	offsetQueryParam := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetQueryParam)
	if err != nil {
		er.SendError(c, er.ErrInvalidValueInOffset)
		return
	}
	videos, err := v.videoService.GetVideos(limit, offset)
	if err != nil {
		er.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, videos)
}

// SearchVideos searches videos from database based on the search string
func (v videoController) SearchVideos(c *gin.Context) {
	var searchRequest model.SearchVideosRequest
	if err := c.BindJSON(&searchRequest); err != nil {
		er.SendError(c, err)
		return
	}
	if searchRequest.SearchString == "" {
		er.SendError(c, er.ErrSearchStringRequired)
		return
	}
	videos, err := v.videoService.SearchVideos(searchRequest.SearchString)
	if err != nil {
		er.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, videos)
}
