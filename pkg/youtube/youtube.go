package youtube

import (
	"context"
	"encoding/json"
	"github.com/Gohelraj/youtube-search-api/api/model"
	"github.com/Gohelraj/youtube-search-api/api/repository"
	"github.com/Gohelraj/youtube-search-api/config"
	"github.com/Gohelraj/youtube-search-api/pkg/ampq"
	"github.com/Gohelraj/youtube-search-api/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/youtube/v3"
)

// SearchVideosFromYoutube searches videos from YouTube and push the videos to queue
func SearchVideosFromYoutube(videoKeyword string, apiKey string, pgxPool *pgxpool.Pool) {
	youtubeRepository := repository.NewYoutubeRepo(pgxPool)
	nextPageToken, err := youtubeRepository.GetAvailableLastPageToken()
	if err != nil {
		log.Printf("Error getting next page token: %v", err)
		return
	}

	service, err := youtube.NewService(context.TODO(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	publishedAfterTime := time.Now().UTC().AddDate(0, -1 , 0).Format(time.RFC3339)
	// Make the API call to YouTube.
	call := service.Search.List([]string{"id,snippet"}).
		Q(videoKeyword).
		PageToken(nextPageToken).
		Order("date").
		Type("video").
		PublishedAfter(publishedAfterTime).
		MaxResults(10)

	response, err := call.Do()
	if err != nil {
		log.Printf("Error making YouTube API call: %v", err)
		return
	}
	if response == nil {
		log.Printf("API call failed: %v", err)
		return
	}

	// Forbidden error is returned when the API key is invalid.
	if response.HTTPStatusCode == http.StatusForbidden {
		// Retry the request with new API key.
		apiKeyIndex := utils.GetIndexOf(apiKey, config.Conf.GoogleAPIKeys)
		if apiKeyIndex+1 < len(config.Conf.GoogleAPIKeys) {
			apiKey = config.Conf.GoogleAPIKeys[apiKeyIndex+1]
			log.Printf("Retrying with new API key")
			SearchVideosFromYoutube(videoKeyword, apiKey, pgxPool)
		} else {
			log.Printf("No more API keys to retry with. Exiting.")
			os.Exit(1)
		}
	}

	var videos []model.VideoMetadata
	for _, item := range response.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Printf("Error parsing publishedAt: %v", err)
		}
		videos = append(videos, model.VideoMetadata{
			ID:           item.Id.VideoId,
			Title:        item.Snippet.Title,
			Description:  &item.Snippet.Description,
			PublishedAt:  publishedAt,
			ThumbnailURL: &item.Snippet.Thumbnails.Medium.Url,
		})
	}

	youtubeVideosQueue := ampq.NewQueue(config.Conf.Ampq.Url, config.Conf.Ampq.QueueName)
	videosData, err := json.Marshal(videos)
	if err != nil {
		log.Printf("Error marshalling videos: %v", err)
		return
	}
	youtubeVideosQueue.Send(videosData)

	// Mark last used page token as used to avoid using it again.
	go youtubeRepository.MarkPageTokenAsUsed(nextPageToken)

	if response.NextPageToken != "" {
		err := youtubeRepository.InsertNextPageToken(response.NextPageToken)
		if err != nil {
			log.Printf("Error inserting next page token: %v", err)
		}
	}
}

// ProcessYoutubeVideosFromQueue processes videos from queue and inserts them into database
func ProcessYoutubeVideosFromQueue(pgxPool *pgxpool.Pool) {
	youtubeVideosQueue := ampq.NewQueue(config.Conf.Ampq.Url, config.Conf.Ampq.QueueName)
	queueMessages, _ := youtubeVideosQueue.Consumer()
	for queueMessage := range queueMessages {
		var videos []model.VideoMetadata
		err := json.Unmarshal(queueMessage.Body, &videos)
		if err != nil {
			log.Printf("Error unmarshalling videos: %v", err)
			continue
		}
		youtubeRepository := repository.NewYoutubeRepo(pgxPool)
		err = youtubeRepository.InsertVideos(videos)
		if err != nil {
			// if error occurred while inserting videos data into database, then retry the request
			_ = queueMessage.Nack(false, true)
			log.Printf("Error inserting videos: %v", err)
			continue
		}
		err = queueMessage.Ack(false)
		if err != nil {
			log.Printf("Error inserting videos: %v", err)
			continue
		}
	}
}
