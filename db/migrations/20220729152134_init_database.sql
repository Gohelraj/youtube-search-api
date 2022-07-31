-- migrate:up
CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    youtube_id VARCHAR(20) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(5000) NULL,
    published_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    thumbnail_url VARCHAR(500) NOT NULL,
    document_with_weights tsvector NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_videos_title_description_index ON videos (title, description);
CREATE INDEX IF NOT EXISTS idx_videos_published_at ON videos (published_at);
CREATE INDEX idx_videos_document_with_weights ON videos USING GIN(document_with_weights);
CREATE UNIQUE INDEX ON videos (youtube_id);

CREATE TABLE IF NOT EXISTS page_tokens (
    id SERIAL PRIMARY KEY,
    next_page_token VARCHAR(20) NOT NULL,
    published_after_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_page_tokens_next_page_token_is_used ON page_tokens (next_page_token, is_used);
CREATE UNIQUE INDEX ON page_tokens (next_page_token);

CREATE FUNCTION videos_tsvector_trigger() RETURNS trigger as $$
BEGIN
    NEW.document_with_weights :=
        setweight(to_tsvector('english', NEW.title), 'A')
        || setweight(to_tsvector('english', NEW.description), 'B');
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER tsvupdate BEFORE INSERT OR UPDATE
    ON videos FOR EACH ROW EXECUTE PROCEDURE videos_tsvector_trigger();

-- migrate:down
DROP TABLE IF EXISTS page_tokens;
DROP INDEX IF EXISTS idx_videos_title_description_index;
DROP INDEX IF EXISTS idx_videos_published_at;
DROP TRIGGER IF EXISTS tsvupdate ON videos;
DROP FUNCTION IF EXISTS videos_tsvector_trigger;
DROP TABLE IF EXISTS videos;
