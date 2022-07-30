package model

import (
	"time"
)

// VideoMetadata video's metadata
type VideoMetadata struct {
	ID           int64     `json:"id,omitempty"`
	YoutubeID    string    `json:"youtubeId"`
	Title        string    `json:"title"`
	Description  *string   `json:"description,omitempty"`
	PublishedAt  time.Time `json:"publishedAt"`
	ThumbnailURL *string   `json:"thumbnailUrl,omitempty"`
}

type SearchVideosRequest struct {
	SearchString string `json:"searchString"`
}
