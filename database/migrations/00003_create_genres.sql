-- +goose Up
-- +goose StatementBegin
CREATE TABLE genres
(
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX genres_name_lower_unique_idx ON genres (LOWER(name));
CREATE INDEX genres_name_trgm_idx ON genres USING GIN (name gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS genres;
-- +goose StatementEnd
