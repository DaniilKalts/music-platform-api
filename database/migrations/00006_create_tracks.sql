-- +goose Up
-- +goose StatementBegin
CREATE TABLE tracks
(
    id               UUID PRIMARY KEY,
    title            TEXT        NOT NULL,
    artist_id        UUID        NOT NULL REFERENCES artists (id) ON DELETE RESTRICT,
    album_id         UUID        NOT NULL REFERENCES albums (id) ON DELETE RESTRICT,
    genre_id         UUID        NOT NULL REFERENCES genres (id) ON DELETE RESTRICT,
    duration_seconds INTEGER     NOT NULL CHECK (duration_seconds > 0),
    file_url         TEXT        NOT NULL,
    deleted_at       TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX tracks_artist_id_idx ON tracks (artist_id);
CREATE INDEX tracks_album_id_idx ON tracks (album_id);
CREATE INDEX tracks_genre_id_idx ON tracks (genre_id);
CREATE INDEX tracks_not_deleted_created_at_idx ON tracks (created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX tracks_not_deleted_title_trgm_idx ON tracks USING GIN (title gin_trgm_ops) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracks;
-- +goose StatementEnd
