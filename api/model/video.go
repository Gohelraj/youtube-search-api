package model

import (
	"time"
)

// VideoMetadata video's metadata
type VideoMetadata struct {
	ID           string    `json:"videoId"`
	Title        string    `json:"title"`
	Description  *string   `json:"description,omitempty"`
	PublishedAt  time.Time `json:"publishedAt"`
	ThumbnailURL *string   `json:"thumbnailUrl,omitempty"`
}
