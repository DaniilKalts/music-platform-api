-- +goose Up
-- +goose StatementBegin
CREATE TABLE artists
(
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE albums
(
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE genres
(
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE UNIQUE INDEX artists_name_lower_unique_idx ON artists (LOWER(name));
CREATE UNIQUE INDEX albums_name_lower_unique_idx ON albums (LOWER(name));
CREATE UNIQUE INDEX genres_name_lower_unique_idx ON genres (LOWER(name));

CREATE INDEX tracks_artist_id_idx ON tracks (artist_id);
CREATE INDEX tracks_album_id_idx ON tracks (album_id);
CREATE INDEX tracks_genre_id_idx ON tracks (genre_id);

CREATE INDEX tracks_not_deleted_created_at_idx ON tracks (created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX tracks_not_deleted_title_trgm_idx ON tracks USING GIN (title gin_trgm_ops) WHERE deleted_at IS NULL;
CREATE INDEX artists_name_trgm_idx ON artists USING GIN (name gin_trgm_ops);
CREATE INDEX albums_name_trgm_idx ON albums USING GIN (name gin_trgm_ops);
CREATE INDEX genres_name_trgm_idx ON genres USING GIN (name gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS albums;
DROP TABLE IF EXISTS artists;
-- +goose StatementEnd
