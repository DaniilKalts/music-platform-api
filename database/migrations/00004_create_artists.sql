-- +goose Up
-- +goose StatementBegin
CREATE TABLE artists
(
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX artists_name_lower_unique_idx ON artists (LOWER(name));
CREATE INDEX artists_name_trgm_idx ON artists USING GIN (name gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS artists;
-- +goose StatementEnd
