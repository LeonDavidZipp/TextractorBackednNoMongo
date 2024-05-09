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

-- name: ListAccounts : many
SELECT * FROM accounts
ORDER BY $1
LIMIT $2
OFFSET $3;

-- name: DeleteAccount :exec
REMOVE FROM accounts
WHERE id = $1;

-- name: AddImage :one

-- name: UpdateAccount :one