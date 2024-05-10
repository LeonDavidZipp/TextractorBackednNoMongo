// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql"
	"time"
)


type Account struct {
	ID         int64          `json:"id"`
	Owner      string         `json:"owner"`
	Email      string         `json:"email"`
	// google calls it sub
	GoogleID   sql.NullString `json:"google_id"`
	// facebook calls it ...
	FacebookID sql.NullString `json:"facebook_id"`
	ImageCount int64          `json:"image_count"`
	Subscribed bool           `json:"subscribed"`
	CreatedAt  time.Time      `json:"created_at"`
}
