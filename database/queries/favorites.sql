-- name: AddFavorite :one
INSERT INTO favorites (user_id, track_id)
SELECT sqlc.arg('user_id'), sqlc.arg('track_id')
WHERE EXISTS (
    SELECT 1
    FROM tracks
    WHERE id = sqlc.arg('track_id')
      AND deleted_at IS NULL
)
RETURNING user_id, track_id, created_at;

-- name: ListFavoritesByUserID :many
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
       f.created_at AS favorited_at
FROM favorites f
JOIN tracks t ON t.id = f.track_id
JOIN artists ar ON ar.id = t.artist_id
JOIN albums al ON al.id = t.album_id
JOIN genres g ON g.id = t.genre_id
WHERE f.user_id = $1
  AND t.deleted_at IS NULL
ORDER BY f.created_at DESC, t.id;

-- name: RemoveFavorite :one
DELETE FROM favorites
WHERE user_id = $1
  AND track_id = $2
RETURNING user_id, track_id;

-- name: CountFavoritesByUserID :one
SELECT COUNT(*)
FROM favorites
WHERE user_id = $1;

-- name: FavoriteExists :one
SELECT EXISTS (
    SELECT 1
    FROM favorites
    WHERE user_id = $1
      AND track_id = $2
);
