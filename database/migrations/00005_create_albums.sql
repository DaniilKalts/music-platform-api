-- +goose Up
-- +goose StatementBegin
CREATE TABLE albums
(
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX albums_name_lower_unique_idx ON albums (LOWER(name));
CREATE INDEX albums_name_trgm_idx ON albums USING GIN (name gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS albums;
-- +goose StatementEnd
