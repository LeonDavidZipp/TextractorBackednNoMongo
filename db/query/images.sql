-- name: CreateImage :one
INSERT INTO images (name) VALUES ($1);

-- name: GetImage :one
SELECT * FROM images
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE;

-- name: GetImageForUpdate :one
SELECT * FROM images
WHERE id = $1
LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
WHERE id = $1;
LIMIT $2
OFFSET $3;

-- name: DeleteImages :exec
DELETE FROM images
WHERE id = ANY($1)

-- name: UpdateImageText :one
UPDATE images
SET text = sqlc.arg(text)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateImageURL :one
UPDATE images
SET url = sqlc.arg(url)
WHERE id = sqlc.arg(id)
RETURNING *;
