package db

import (
	"context"
	"database/sql"
	"fmt"
	mongodb "go.mongodb.org/mongo-driver/mongo"
)


type Store struct {
	Querier
	UserDB *sql.DB
	ImageDBClient *mongodb.Client
}

type UploadImageTransactionParams struct {
	Filepath string `json:"filepath"`
	AccountID int64 `json:"account_id"`
}

type UploadImageTransactionResult struct {

}

// Upload Image handles uploading the necessary data and image to the databases.
func (store *Store) UploadImage(
	ctx context.Context,
	arg UploadImageTransactionParams
) UploadImageTransactionResult {
	err := store.execTransaction(ctx, arg.Filepath, UploadToPostgres, UploadToMongo)
}

// Uploads the data to postgres accounts table
func (store *Store) uploadToPostgres(querie *Queries) error {
	var result UploadImageTransactionResult


}

// Uploads the image to mongodb table
func (store *Store) uploadToMongo(mongoCtx ) error {
	mongoDB, err := ImageDBClient.Database()
}

// execTransaction creates a "rollback-able" transaction.
// It takes as arguments
// 	the context
// 	the filepath to the file to upload
// 	the transcribed text
// 	a function uploading the account data to accounts postgresql table
// 	a function uploading image to mongodb table
func (store *Store) execTransaction(
	ctx context.Context,
	filepath string,
	text string,
	fnSql func(*Queries) error,
	fnMongo func(mongo.SessionContext, string, string) error
	) error {
	// postgres
	sqlTransaction, err := store.UserDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sqlQuerie := New(sqlTransaction)
	err = fnSql(sqlQuerie)
	if err != nil {
		if rollbackErr := sqlTransaction.Rollback(), rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}

	// mongodb
	session, err := store.ImageDBClient.StartSession()
	if err != nil {
		if rollbackErr := sqlTransaction.Rollback(), rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}
	defer session.EndSession()

	err = session.StartTransaction()
	if err != nil {
		if rollbackErr := sqlTransaction.Rollback(), rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}

	// TODO: move database into fnMongo, implement uploading image in fnMongo
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		return fnMongo(sessionContext, filePath, text) // Pass file path to fnMongo
	})
	if er != nil {
		if rollbackErr := sqlTransaction.Rollback(), rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}

	// commit
	err = sqlTransaction.Commit()
	if err != nil {
		return err
	}
	err = session.CommitTransaction(ctx)
	if err != nil {
		return err
	}
	return nil 
}
