package repository

import (
	"context"
	"github.com/Gohelraj/youtube-search-api/api/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type YoutubeRepository interface {
	InsertVideos(videos []model.VideoMetadata) error
	InsertNextPageToken(pageToken string) error
	GetAvailableLastPageToken() (pageToken string, err error)
	MarkPageTokenAsUsed(pageToken string) error
	GetVideos(limit int, offset int) ([]model.VideoMetadata, error)
	SearchVideos(searchString string) ([]model.VideoMetadata, error)
}

type youtubeRepository struct {
	pgxPool *pgxpool.Pool
}

func NewYoutubeRepo(pgxPool *pgxpool.Pool) YoutubeRepository {
	return youtubeRepository{
		pgxPool: pgxPool,
	}
}

func (youtubeRepo youtubeRepository) InsertVideos(videos []model.VideoMetadata) error {
	conn, err := youtubeRepo.pgxPool.Acquire(context.Background())
	if err != nil {
		log.Printf("Couldn't get a connection with the database. Reason %v", err)
		return err
	}
	defer conn.Release()
	batch := &pgx.Batch{}
	for _, video := range videos {
		currentTime := time.Now().UTC()
		// here we are using "ON CONFLICT DO NOTHING" to avoid duplicate entries/duplicate primary key violation errors
		batch.Queue("INSERT INTO videos (youtube_id, title, description, published_at, created_at, updated_at, thumbnail_url) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING", video.YoutubeID, video.Title, video.Description, video.PublishedAt, currentTime, currentTime, video.ThumbnailURL)
	}
	result := conn.SendBatch(context.Background(), batch)
	for i := 0; i < batch.Len(); i++ {
		_, err := result.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (youtubeRepo youtubeRepository) GetAvailableLastPageToken() (pageToken string, err error) {
	conn, err := youtubeRepo.pgxPool.Acquire(context.Background())
	if err != nil {
		log.Printf("Couldn't get a connection with the database. Reason %v", err)
		return "", err
	}
	defer conn.Release()
	row := conn.QueryRow(context.Background(), "SELECT next_page_token FROM page_tokens WHERE is_used = false ORDER BY created_at DESC LIMIT 1")
	err = row.Scan(&pageToken)
	if err != nil && err != pgx.ErrNoRows {
		return "", err
	}
	return pageToken, nil
}

func (youtubeRepo youtubeRepository) InsertNextPageToken(pageToken string) error {
	conn, err := youtubeRepo.pgxPool.Acquire(context.Background())
	if err != nil {
		log.Printf("Couldn't get a connection with the database. Reason %v", err)
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), "INSERT INTO page_tokens (next_page_token, created_at) VALUES ($1, $2) ON CONFLICT DO NOTHING", pageToken, time.Now().UTC())
	return err
}

func (youtubeRepo youtubeRepository) MarkPageTokenAsUsed(pageToken string) error {
	conn, err := youtubeRepo.pgxPool.Acquire(context.Background())
	if err != nil {
		log.Printf("Couldn't get a connection with the database. Reason %v", err)
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), "UPDATE page_tokens SET is_used = true WHERE next_page_token = $1", pageToken)
	return err
}

func (youtubeRepo youtubeRepository) GetVideos(limit int, offset int) ([]model.VideoMetadata, error) {
	conn, err := youtubeRepo.pgxPool.Acquire(context.Background())
	if err != nil {
		log.Printf("Couldn't get a connection with the database. Reason %v", err)
		return []model.VideoMetadata{}, err
	}
	defer conn.Release()
	videos := []model.VideoMetadata{}
	rows, err := conn.Query(context.Background(), "SELECT id, youtube_id, title, description, published_at, thumbnail_url FROM videos ORDER BY published_at DESC LIMIT $1 OFFSET $2", limit, offset)
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

func (youtubeRepo youtubeRepository) SearchVideos(searchString string) ([]model.VideoMetadata, error) {
	conn, err := youtubeRepo.pgxPool.Acquire(context.Background())
	if err != nil {
		log.Printf("Couldn't get a connection with the database. Reason %v", err)
		return []model.VideoMetadata{}, err
	}
	defer conn.Release()
	videos := []model.VideoMetadata{}
	rows, err := conn.Query(context.Background(), "SELECT id, youtube_id, title, description, published_at, thumbnail_url FROM videos WHERE document_with_weights @@ plainto_tsquery($1) ORDER BY ts_rank(document_with_weights, plainto_tsquery($1)) DESC, published_at DESC", searchString)
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