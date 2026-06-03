-- +goose Up
-- +goose StatementBegin
CREATE TABLE favorites
(
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    track_id   UUID        NOT NULL REFERENCES tracks (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, track_id)
);

CREATE INDEX favorites_track_id_idx ON favorites (track_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS favorites;
-- +goose StatementEnd
