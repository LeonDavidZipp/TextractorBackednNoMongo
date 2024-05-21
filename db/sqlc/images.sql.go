// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: images.sql

package db

import (
	"context"

	"github.com/lib/pq"
)

const createImage = `-- name: CreateImage :one
INSERT INTO images 
(user_id, url, text)
VALUES ($1, $2, $3)
RETURNING id, user_id, url, preview_url, text, created_at
`

type CreateImageParams struct {
	UserID int64  `json:"user_id"`
	Url    string `json:"url"`
	Text   string `json:"text"`
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, createImage, arg.UserID, arg.Url, arg.Text)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.PreviewUrl,
		&i.Text,
		&i.CreatedAt,
	)
	return i, err
}

const deleteImages = `-- name: DeleteImages :exec
DELETE FROM images
WHERE id = ANY($1::bigint[])
`

func (q *Queries) DeleteImages(ctx context.Context, ids []int64) error {
	_, err := q.db.ExecContext(ctx, deleteImages, pq.Array(ids))
	return err
}

const getImageForUpdate = `-- name: GetImageForUpdate :one
SELECT id, user_id, url, preview_url, text, created_at FROM images
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetImageForUpdate(ctx context.Context, id int64) (Image, error) {
	row := q.db.QueryRowContext(ctx, getImageForUpdate, id)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.PreviewUrl,
		&i.Text,
		&i.CreatedAt,
	)
	return i, err
}

const getImageFromSQL = `-- name: GetImageFromSQL :one
SELECT id, user_id, url, preview_url, text, created_at FROM images
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetImageFromSQL(ctx context.Context, id int64) (Image, error) {
	row := q.db.QueryRowContext(ctx, getImageFromSQL, id)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.PreviewUrl,
		&i.Text,
		&i.CreatedAt,
	)
	return i, err
}

const listImages = `-- name: ListImages :many
SELECT id, user_id, url, preview_url, text, created_at FROM images
WHERE user_id = $1
LIMIT $2
OFFSET $3
`

type ListImagesParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListImages(ctx context.Context, arg ListImagesParams) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, listImages, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Url,
			&i.PreviewUrl,
			&i.Text,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateImageText = `-- name: UpdateImageText :one
UPDATE images
SET text = $1
WHERE id = $2
RETURNING id, user_id, url, preview_url, text, created_at
`

type UpdateImageTextParams struct {
	Text string `json:"text"`
	ID   int64  `json:"id"`
}

func (q *Queries) UpdateImageText(ctx context.Context, arg UpdateImageTextParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, updateImageText, arg.Text, arg.ID)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.PreviewUrl,
		&i.Text,
		&i.CreatedAt,
	)
	return i, err
}

const updateImageUrl = `-- name: UpdateImageUrl :one
UPDATE images
SET url = $1
WHERE id = $2
RETURNING id, user_id, url, preview_url, text, created_at
`

type UpdateImageUrlParams struct {
	Url string `json:"url"`
	ID  int64  `json:"id"`
}

func (q *Queries) UpdateImageUrl(ctx context.Context, arg UpdateImageUrlParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, updateImageUrl, arg.Url, arg.ID)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Url,
		&i.PreviewUrl,
		&i.Text,
		&i.CreatedAt,
	)
	return i, err
}
