package service

import "github.com/Gohelraj/youtube-search-api/api/repository"

type YoutubeService interface {
}

type youtubeService struct {
	youtubeRepository repository.YoutubeRepository
}

func NewYoutubeService(r repository.YoutubeRepository) YoutubeService {
	return youtubeService{
		youtubeRepository: r,
	}
}
