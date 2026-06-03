-- name: CreateArtist :one
INSERT INTO artists (id, name)
VALUES ($1, $2)
RETURNING id, name, created_at, updated_at;

-- name: GetArtistByID :one
SELECT id, name, created_at, updated_at
FROM artists
WHERE id = $1;

-- name: GetArtistByName :one
SELECT id, name, created_at, updated_at
FROM artists
WHERE LOWER(name) = LOWER($1);

-- name: FindOrCreateArtist :one
INSERT INTO artists (id, name)
VALUES ($1, $2)
ON CONFLICT (LOWER(name)) DO UPDATE
SET name = artists.name
RETURNING id, name, created_at, updated_at;

-- name: CreateAlbum :one
INSERT INTO albums (id, name)
VALUES ($1, $2)
RETURNING id, name, created_at, updated_at;

-- name: GetAlbumByID :one
SELECT id, name, created_at, updated_at
FROM albums
WHERE id = $1;

-- name: GetAlbumByName :one
SELECT id, name, created_at, updated_at
FROM albums
WHERE LOWER(name) = LOWER($1);

-- name: FindOrCreateAlbum :one
INSERT INTO albums (id, name)
VALUES ($1, $2)
ON CONFLICT (LOWER(name)) DO UPDATE
SET name = albums.name
RETURNING id, name, created_at, updated_at;

-- name: CreateGenre :one
INSERT INTO genres (id, name)
VALUES ($1, $2)
RETURNING id, name, created_at, updated_at;

-- name: ListGenres :many
SELECT id, name, created_at, updated_at
FROM genres
ORDER BY name;

-- name: GetGenreByID :one
SELECT id, name, created_at, updated_at
FROM genres
WHERE id = $1;

-- name: GetGenreByName :one
SELECT id, name, created_at, updated_at
FROM genres
WHERE LOWER(name) = LOWER($1);
