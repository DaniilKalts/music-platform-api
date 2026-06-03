-- name: CreatePlaylist :one
INSERT INTO playlists (id, user_id, name, description)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, name, description, created_at, updated_at;

-- name: ListPlaylistsByUserID :many
SELECT id, user_id, name, description, created_at, updated_at
FROM playlists
WHERE user_id = $1
ORDER BY created_at DESC, id;

-- name: GetPlaylistByIDForUser :one
SELECT id, user_id, name, description, created_at, updated_at
FROM playlists
WHERE id = $1
  AND user_id = $2;

-- name: UpdatePlaylist :one
UPDATE playlists
SET name = $3,
    description = $4,
    updated_at = NOW()
WHERE id = $1
  AND user_id = $2
RETURNING id, user_id, name, description, created_at, updated_at;

-- name: DeletePlaylist :one
DELETE FROM playlists
WHERE id = $1
  AND user_id = $2
RETURNING id;

-- name: CountPlaylistsByUserID :one
SELECT COUNT(*)
FROM playlists
WHERE user_id = $1;

-- name: AddTrackToPlaylist :one
INSERT INTO playlist_tracks (playlist_id, track_id)
SELECT sqlc.arg('playlist_id'), sqlc.arg('track_id')
WHERE EXISTS (
    SELECT 1
    FROM playlists
    WHERE id = sqlc.arg('playlist_id')
      AND user_id = sqlc.arg('user_id')
)
  AND EXISTS (
    SELECT 1
    FROM tracks
    WHERE id = sqlc.arg('track_id')
      AND deleted_at IS NULL
)
RETURNING playlist_id, track_id, added_at;

-- name: RemoveTrackFromPlaylist :one
DELETE FROM playlist_tracks
WHERE playlist_id = sqlc.arg('playlist_id')
  AND track_id = sqlc.arg('track_id')
  AND EXISTS (
      SELECT 1
      FROM playlists
      WHERE id = sqlc.arg('playlist_id')
        AND user_id = sqlc.arg('user_id')
  )
RETURNING playlist_id, track_id;

-- name: ListPlaylistTracks :many
SELECT t.id,
       t.title,
       t.artist_id,
       ar.name AS artist_name,
       t.album_id,
       al.name AS album_name,
       t.genre_id,
       g.name AS genre_name,
       t.duration_seconds,
       t.file_url,
       t.deleted_at,
       t.created_at,
       t.updated_at,
       pt.added_at
FROM playlist_tracks pt
JOIN playlists p ON p.id = pt.playlist_id
JOIN tracks t ON t.id = pt.track_id
JOIN artists ar ON ar.id = t.artist_id
JOIN albums al ON al.id = t.album_id
JOIN genres g ON g.id = t.genre_id
WHERE pt.playlist_id = sqlc.arg('playlist_id')
  AND p.user_id = sqlc.arg('user_id')
  AND t.deleted_at IS NULL
ORDER BY pt.added_at DESC, t.id;
