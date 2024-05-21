-- name: CreateImage :one
INSERT INTO images 
(user_id, url, preview_url, text)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetImageFromSQL :one
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
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: DeleteImages :exec
DELETE FROM images
WHERE id = ANY(sqlc.arg(ids)::bigint[]);

-- name: UpdateImageText :one
UPDATE images
SET text = sqlc.arg(text)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateImageUrl :one
UPDATE images
SET url = sqlc.arg(url)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateImagePreviewUrl :one
UPDATE images
SET preview_url = sqlc.arg(url)
WHERE id = sqlc.arg(id)
RETURNING *;
