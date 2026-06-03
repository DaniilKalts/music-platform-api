-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlist_tracks
(
    playlist_id UUID        NOT NULL REFERENCES playlists (id) ON DELETE CASCADE,
    track_id    UUID        NOT NULL REFERENCES tracks (id) ON DELETE RESTRICT,
    added_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (playlist_id, track_id)
);

CREATE INDEX playlist_tracks_track_id_idx ON playlist_tracks (track_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS playlist_tracks;
-- +goose StatementEnd
