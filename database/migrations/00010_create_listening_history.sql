-- +goose Up
-- +goose StatementBegin
CREATE TABLE listening_history
(
    id          UUID PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    track_id    UUID        NOT NULL REFERENCES tracks (id) ON DELETE RESTRICT,
    listened_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX listening_history_user_id_listened_at_idx ON listening_history (user_id, listened_at DESC);
CREATE INDEX listening_history_track_id_idx ON listening_history (track_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS listening_history;
-- +goose StatementEnd
