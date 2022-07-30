-- migrate:up
CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    youtube_id VARCHAR(20) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(5000) NULL,
    published_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    thumbnail_url VARCHAR(500) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_videos_title_description_index ON videos (title, description);
CREATE INDEX IF NOT EXISTS idx_videos_published_at ON videos (published_at);
CREATE UNIQUE INDEX ON videos (youtube_id);

CREATE TABLE IF NOT EXISTS page_tokens (
    id SERIAL PRIMARY KEY,
    next_page_token VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_page_tokens_next_page_token_is_used ON page_tokens (next_page_token, is_used);
CREATE UNIQUE INDEX ON page_tokens (next_page_token);

-- migrate:down
DROP TABLE IF EXISTS page_tokens;
DROP INDEX IF EXISTS idx_videos_title_description_index;
DROP INDEX IF EXISTS idx_videos_published_at;
DROP TABLE IF EXISTS videos;
