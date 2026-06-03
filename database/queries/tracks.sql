-- name: CreateTrack :one
WITH inserted AS (
    INSERT INTO tracks (id, title, artist_id, album_id, genre_id, duration_seconds, file_url)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, title, artist_id, album_id, genre_id, duration_seconds, file_url, deleted_at, created_at, updated_at
)
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
       t.updated_at
FROM inserted t
JOIN artists ar ON ar.id = t.artist_id
JOIN albums al ON al.id = t.album_id
JOIN genres g ON g.id = t.genre_id;

-- name: GetTrackByID :one
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
       t.updated_at
FROM tracks t
JOIN artists ar ON ar.id = t.artist_id
JOIN albums al ON al.id = t.album_id
JOIN genres g ON g.id = t.genre_id
WHERE t.id = $1
  AND t.deleted_at IS NULL;

-- name: ListTracks :many
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
       t.updated_at
FROM tracks t
JOIN artists ar ON ar.id = t.artist_id
JOIN albums al ON al.id = t.album_id
JOIN genres g ON g.id = t.genre_id
WHERE t.deleted_at IS NULL
ORDER BY t.created_at DESC, t.id
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: SearchTracks :many
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
       t.updated_at
FROM tracks t
JOIN artists ar ON ar.id = t.artist_id
JOIN albums al ON al.id = t.album_id
JOIN genres g ON g.id = t.genre_id
WHERE t.deleted_at IS NULL
  AND (
      t.title ILIKE '%' || sqlc.arg('query')::text || '%'
      OR ar.name ILIKE '%' || sqlc.arg('query')::text || '%'
      OR al.name ILIKE '%' || sqlc.arg('query')::text || '%'
      OR g.name ILIKE '%' || sqlc.arg('query')::text || '%'
  )
ORDER BY t.created_at DESC, t.id
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateTrack :one
WITH updated AS (
    UPDATE tracks
    SET title = $2,
        artist_id = $3,
        album_id = $4,
        genre_id = $5,
        duration_seconds = $6,
        file_url = $7,
        updated_at = NOW()
    WHERE tracks.id = $1
      AND tracks.deleted_at IS NULL
    RETURNING tracks.id, tracks.title, tracks.artist_id, tracks.album_id, tracks.genre_id, tracks.duration_seconds, tracks.file_url, tracks.deleted_at, tracks.created_at, tracks.updated_at
)
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
       t.updated_at
FROM updated t
JOIN artists ar ON ar.id = t.artist_id
JOIN albums al ON al.id = t.album_id
JOIN genres g ON g.id = t.genre_id;

-- name: SoftDeleteTrack :one
UPDATE tracks
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE tracks.id = $1
  AND tracks.deleted_at IS NULL
RETURNING tracks.id;

-- name: TrackExists :one
SELECT EXISTS (
    SELECT 1
    FROM tracks
    WHERE id = $1
      AND deleted_at IS NULL
);
