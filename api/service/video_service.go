package service

import (
	"github.com/Gohelraj/youtube-search-api/api/model"
	"github.com/Gohelraj/youtube-search-api/api/repository"
)

type VideoService interface {
	GetVideos(limit int, offset int) ([]model.VideoMetadata, error)
	SearchVideos(searchString string) ([]model.VideoMetadata, error)
}

type videoService struct {
	videoRepository repository.VideoRepository
}

func NewVideoService(r repository.VideoRepository) VideoService {
	return videoService{
		videoRepository: r,
	}
}

func (v videoService) GetVideos(limit int, offset int) ([]model.VideoMetadata, error) {
	return v.videoRepository.GetVideos(limit, offset)
}

func (v videoService) SearchVideos(searchString string) ([]model.VideoMetadata, error) {
	return v.videoRepository.SearchVideos(searchString)
}
