package repository

import (
	"context"
	"github.com/Gohelraj/youtube-search-api/api/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type VideoRepository interface {
	InsertVideos(videos []model.VideoMetadata) error
	InsertNextPageToken(pageToken string, publishedAfterDateTime time.Time) error
	GetAvailableLastPageToken() (pageToken string, publishedAfterDateTime time.Time, err error)
	MarkPageTokenAsUsed(pageToken string) error
	GetVideos(limit int, offset int) ([]model.VideoMetadata, error)
	SearchVideos(searchString string) ([]model.VideoMetadata, error)
	GetLastPublishedAtDateTime() (time.Time, error)
}

type videoRepository struct {
	pgxPool *pgxpool.Pool
}

func NewVideoRepo(pgxPool *pgxpool.Pool) VideoRepository {
	return videoRepository{
		pgxPool: pgxPool,
	}
}

// InsertVideos batch inserts videos into the database.
func (videoRepo videoRepository) InsertVideos(videos []model.VideoMetadata) error {
	batch := &pgx.Batch{}
	for _, video := range videos {
		currentTime := time.Now().UTC()
		// here we are using "ON CONFLICT DO NOTHING" to avoid duplicate entries/duplicate primary key violation errors
		batch.Queue("INSERT INTO videos (youtube_id, title, description, published_at, created_at, updated_at, thumbnail_url) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING", video.YoutubeID, video.Title, video.Description, video.PublishedAt, currentTime, currentTime, video.ThumbnailURL)
	}
	result := videoRepo.pgxPool.SendBatch(context.Background(), batch)
	for i := 0; i < batch.Len(); i++ {
		_, err := result.Exec()
		if err != nil {
			return err
		}
	}
	defer result.Close()
	return nil
}

// GetAvailableLastPageToken returns the next page token that is not used.
func (videoRepo videoRepository) GetAvailableLastPageToken() (pageToken string, publishedAfterDateTime time.Time, err error) {
	row := videoRepo.pgxPool.QueryRow(context.Background(), "SELECT next_page_token, published_after_time FROM page_tokens WHERE is_used = false ORDER BY created_at DESC LIMIT 1")
	err = row.Scan(&pageToken, &publishedAfterDateTime)
	if err != nil && err != pgx.ErrNoRows {
		return "", time.Time{}, err
	}
	return pageToken, publishedAfterDateTime, nil
}

// InsertNextPageToken inserts the next page token into the database.
func (videoRepo videoRepository) InsertNextPageToken(pageToken string, publishedAfterDateTime time.Time) error {
	_, err := videoRepo.pgxPool.Exec(context.Background(), "INSERT INTO page_tokens (next_page_token, published_after_time, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", pageToken, publishedAfterDateTime.Format(time.RFC3339), time.Now().UTC())
	return err
}

// MarkPageTokenAsUsed marks the page token as used.
func (videoRepo videoRepository) MarkPageTokenAsUsed(pageToken string) error {
	_, err := videoRepo.pgxPool.Exec(context.Background(), "UPDATE page_tokens SET is_used = true WHERE next_page_token = $1", pageToken)
	return err
}

// GetVideos returns the videos from the database.
func (videoRepo videoRepository) GetVideos(limit int, offset int) ([]model.VideoMetadata, error) {
	videos := []model.VideoMetadata{}
	rows, err := videoRepo.pgxPool.Query(context.Background(), "SELECT id, youtube_id, title, description, published_at, thumbnail_url FROM videos ORDER BY published_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var video model.VideoMetadata
		err = rows.Scan(&video.ID, &video.YoutubeID, &video.Title, &video.Description, &video.PublishedAt, &video.ThumbnailURL)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

// SearchVideos search videos from the database using full text search based on given search string.
func (videoRepo videoRepository) SearchVideos(searchString string) ([]model.VideoMetadata, error) {
	videos := []model.VideoMetadata{}
	rows, err := videoRepo.pgxPool.Query(context.Background(), "SELECT id, youtube_id, title, description, published_at, thumbnail_url FROM videos WHERE document_with_weights @@ plainto_tsquery($1) ORDER BY ts_rank(document_with_weights, plainto_tsquery($1)) DESC, published_at DESC", searchString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var video model.VideoMetadata
		err = rows.Scan(&video.ID, &video.YoutubeID, &video.Title, &video.Description, &video.PublishedAt, &video.ThumbnailURL)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

// GetLastPublishedAtDateTime returns the latest published at date time from the database.
func (videoRepo videoRepository) GetLastPublishedAtDateTime() (time.Time, error) {
	row := videoRepo.pgxPool.QueryRow(context.Background(), "SELECT published_at FROM videos ORDER BY published_at DESC LIMIT 1")
	var publishedAt time.Time
	err := row.Scan(&publishedAt)
	if err != nil && err != pgx.ErrNoRows {
		return publishedAt, err
	}
	return publishedAt, nil
}
