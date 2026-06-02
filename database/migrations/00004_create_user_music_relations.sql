-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlists
(
    id          UUID PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE playlist_tracks
(
    playlist_id UUID        NOT NULL REFERENCES playlists (id) ON DELETE CASCADE,
    track_id    UUID        NOT NULL REFERENCES tracks (id) ON DELETE RESTRICT,
    added_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (playlist_id, track_id)
);

CREATE TABLE favorites
(
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    track_id   UUID        NOT NULL REFERENCES tracks (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, track_id)
);

CREATE TABLE listening_history
(
    id          UUID PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    track_id    UUID        NOT NULL REFERENCES tracks (id) ON DELETE RESTRICT,
    listened_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX playlists_user_id_name_lower_unique_idx ON playlists (user_id, LOWER(name));

CREATE INDEX playlist_tracks_track_id_idx ON playlist_tracks (track_id);
CREATE INDEX favorites_track_id_idx ON favorites (track_id);
CREATE INDEX listening_history_user_id_listened_at_idx ON listening_history (user_id, listened_at DESC);
CREATE INDEX listening_history_track_id_idx ON listening_history (track_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS listening_history;
DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS playlist_tracks;
DROP TABLE IF EXISTS playlists;
-- +goose StatementEnd
