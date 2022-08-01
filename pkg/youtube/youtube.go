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
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/youtube/v3"
)

// SearchVideosFromYoutubeAndAddToQueue searches videos from YouTube and push the videos to queue
func SearchVideosFromYoutubeAndAddToQueue(videoKeyword string, pgxPool *pgxpool.Pool) {
	youtubeRepository := repository.NewVideoRepo(pgxPool)
	nextPageToken, publishedAfter, err := youtubeRepository.GetAvailableLastPageToken()
	if err != nil {
		log.Printf("Error getting next page token: %v", err)
		return
	}
	if nextPageToken == "" {
		publishedAfter, err = youtubeRepository.GetLastPublishedAtDateTime()
		if err != nil {
			log.Printf("Error getting last published at date time: %v", err)
			return
		}
		if publishedAfter.IsZero() {
			publishedAfter = time.Now().UTC().Add(-2 * time.Hour)
		}
	}
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(config.Conf.ActiveGoogleAPIKey))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
		return
	}
	// Prepare the API call.
	call := service.Search.List([]string{"id,snippet"}).
		Q(videoKeyword).
		PageToken(nextPageToken).
		Order("date").
		Type("video").
		PublishedAfter(publishedAfter.Format(time.RFC3339)).
		MaxResults(50)

	// Make the API call to YouTube.
	response, err := call.Do()
	if err != nil {
		if apiError, ok := err.(*googleapi.Error); ok {
			// Forbidden error is returned when the API key quota exhausted.
			if apiError.Code == http.StatusForbidden {
				retryWithNewAPIKeyWhenForbidden(videoKeyword, pgxPool)
			}
		}
		log.Printf("Error making YouTube API call: %v", err)
		return
	}
	if response == nil {
		log.Printf("API call failed: %v", err)
		return
	}

	if response.HTTPStatusCode == http.StatusForbidden {
		retryWithNewAPIKeyWhenForbidden(videoKeyword, pgxPool)
	}

	var videos []model.VideoMetadata
	for _, item := range response.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Printf("Error parsing publishedAt: %v", err)
			continue
		}
		videos = append(videos, model.VideoMetadata{
			YoutubeID:    item.Id.VideoId,
			Title:        item.Snippet.Title,
			Description:  &item.Snippet.Description,
			PublishedAt:  publishedAt,
			ThumbnailURL: &item.Snippet.Thumbnails.Medium.Url,
		})
	}

	if len(videos) > 0 {
		youtubeVideosQueue := ampq.NewQueue(config.Conf.Ampq.Url, config.Conf.Ampq.QueueName)
		videosData, err := json.Marshal(videos)
		if err != nil {
			log.Printf("Error marshalling videos: %v", err)
			return
		}
		youtubeVideosQueue.Send(videosData)
		log.Println("Queued youtube videos")
	}

	// Mark last used page token as used to avoid using it again.
	go youtubeRepository.MarkPageTokenAsUsed(nextPageToken)

	if response.NextPageToken != "" {
		// Store next page token to be used in next search.
		err := youtubeRepository.InsertNextPageToken(response.NextPageToken, publishedAfter)
		if err != nil {
			log.Printf("Error inserting next page token: %v", err)
			return
		}
	}
}

// ProcessYoutubeVideosFromQueue processes videos from queue and inserts them into database
func ProcessYoutubeVideosFromQueue(pgxPool *pgxpool.Pool) {
	youtubeVideosQueue := ampq.NewQueue(config.Conf.Ampq.Url, config.Conf.Ampq.QueueName)
	queueMessages, err := youtubeVideosQueue.Consumer()
	if err != nil {
		log.Printf("Error getting queue messages: %v", err)
		return
	}
	for queueMessage := range queueMessages {
		var videos []model.VideoMetadata
		err := json.Unmarshal(queueMessage.Body, &videos)
		if err != nil {
			log.Printf("Error unmarshalling videos: %v", err)
			continue
		}
		youtubeRepository := repository.NewVideoRepo(pgxPool)
		err = youtubeRepository.InsertVideos(videos)
		if err != nil {
			// if error occurred while inserting videos data into database, then retry the request
			_ = queueMessage.Nack(false, true)
			log.Printf("Error inserting videos: %v", err)
			continue
		}
		// if no error occurred while inserting videos data into database, then ack the message
		err = queueMessage.Ack(false)
		if err != nil {
			log.Printf("Error inserting videos: %v", err)
			continue
		}
	}
}

// retryWithNewAPIKeyWhenForbidden retries the request with new API key when forbidden error is returned
func retryWithNewAPIKeyWhenForbidden(videoKeyword string, pgxPool *pgxpool.Pool) {
	// Retry the request with new API key.
	apiKeyIndex := utils.GetIndexOf(config.Conf.ActiveGoogleAPIKey, config.Conf.GoogleAPIKeys)
	if apiKeyIndex+1 < len(config.Conf.GoogleAPIKeys) {
		config.Conf.ActiveGoogleAPIKey = config.Conf.GoogleAPIKeys[apiKeyIndex+1]
		log.Printf("Retrying with new API key")
		SearchVideosFromYoutubeAndAddToQueue(videoKeyword, pgxPool)
	} else {
		log.Printf("Retrying with rotating API keys")
		// If all API keys are used, then retry the request with first API key.
		// This will ensure that the API key rotation will continue.
		config.Conf.ActiveGoogleAPIKey = config.Conf.GoogleAPIKeys[0]
		SearchVideosFromYoutubeAndAddToQueue(videoKeyword, pgxPool)
	}
}
