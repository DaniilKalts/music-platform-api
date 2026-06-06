-- name: CreateListeningHistory :one
INSERT INTO listening_history (id, user_id, track_id)
SELECT sqlc.arg('id'), sqlc.arg('user_id'), sqlc.arg('track_id')
WHERE EXISTS (
    SELECT 1
    FROM tracks
    WHERE id = sqlc.arg('track_id')
      AND deleted_at IS NULL
)
RETURNING id, user_id, track_id, listened_at;

-- name: ListListeningHistoryByUserID :many
SELECT h.id,
       h.user_id,
       t.id AS track_id,
       t.title,
       ar.name AS artist_name,
       h.listened_at
FROM (
    SELECT DISTINCT ON (track_id) id, user_id, track_id, listened_at
    FROM listening_history
    WHERE user_id = $1
    ORDER BY track_id, listened_at DESC, id DESC
) h
JOIN tracks t ON t.id = h.track_id
JOIN artists ar ON ar.id = t.artist_id
WHERE t.deleted_at IS NULL
ORDER BY h.listened_at DESC, h.id
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
