-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    email,
    google_id,
    facebook_id
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;

-- name: UpdateImageCount :one
UPDATE accounts
SET image_count = image_count + sqlc.arg(amount)
RETURNING *;

-- name: UpdateEmail :one
UPDATE accounts
SET email = $2
WHERE id = $1
RETURNING *;

-- name: UpdateSubscribed :one
UPDATE accounts
SET subscribed = $2
WHERE id = $1
RETURNING *;
