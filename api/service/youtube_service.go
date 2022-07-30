package service

import (
	"github.com/Gohelraj/youtube-search-api/api/model"
	"github.com/Gohelraj/youtube-search-api/api/repository"
)

type YoutubeService interface {
	GetVideos(limit int, offset int) ([]model.VideoMetadata, error)
	SearchVideos(searchString string) ([]model.VideoMetadata, error)
}

type youtubeService struct {
	youtubeRepository repository.YoutubeRepository
}

func NewYoutubeService(r repository.YoutubeRepository) YoutubeService {
	return youtubeService{
		youtubeRepository: r,
	}
}

func (y youtubeService) GetVideos(limit int, offset int) ([]model.VideoMetadata, error) {
	return y.youtubeRepository.GetVideos(limit, offset)
}

func (y youtubeService) SearchVideos(searchString string) ([]model.VideoMetadata, error) {
	return y.youtubeRepository.SearchVideos(searchString)
}