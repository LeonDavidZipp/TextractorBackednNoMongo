package db

import (
	"context"
	"database/sql"
	"fmt"
	// mongodb "go.mongodb.org/mongo-driver/mongo"
)


type Store struct {
	Querier
	UserDB *sql.DB
	// ImageDBClient *mongodb.Client
}

func (store *Store) execTransaction(ctx context.Context fn func(*Queries) error) error {
	sql_tx, err := store.UserDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sql_querie := New(sql_tx)
	err = fn
}