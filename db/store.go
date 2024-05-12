package db

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"go.mongodb.org/mongo-driver/mongo"
)


type Store interface {
	db.Querier
	mongodb.MongoQuerier
	UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) UploadImageTransactionResult
}

type SQLMongoStore struct {
	*db.Queries
	*mongodb.MongoQueries
	UserDB *sql.DB
	ImageDB *mongo.Database
}

func NewStore(userDB *sql.DB, imageDB *mongo.Database) Store {
	return &SQLMongoStore{
		Queries: db.New(userDB),
		MongoQueries: mongodb.NewMongo(imageDB),
		UserDB: userDB,
		ImageDB: imageDB,
	}
}

type UploadImageTransactionParams struct {
	AccountID int64  `json:"account_id"`
	Filepath  string `json:"filepath"`
	Text      string `json:"text"`
	Link      string `json:"link"`
}

type UploadImageTransactionResult struct {
	Image    mongodb.Image `json:"image"`
	Uploader db.Account    `json:"uploader"`
}

// Upload Image handles uploading the necessary data and image to the databases.
func (store *SQLMongoStore) UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) UploadImageTransactionResult {
	var result UploadImageTransactionResult
	
	err := store.execTransaction(
		ctx,
		arg.Filepath,
		func(q *Queries) error {
			return uploadToPostgres(arg.AccountID, &result),
			},
		func(q *Queries) error {
			return uploadToMongo(arg, &result),
		},
	)
	if err != nil {
		return UploadImageTransactionResult{}
	}
}

// Uploads the data to postgres accounts table
func (store *Store) uploadToPostgres(AccountID int64, result *UploadImageTransactionResult) error {
	result.Uploader, err := store.UpdateImageCount(ctx, UpdateImageCountParams{
		Amount: 1,
		ID: arg.AccountID,
	})
	if err != nil {
		return err
	}
	return nil
}

// Uploads the image to mongodb table
func (store *Store) uploadToMongo(arg UploadImageTransactionParams, result *UploadImageTransactionResult) error {
	mongoDB, err := store.ImageDBClient.Database()
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
	arg UploadImageTransactionParams,
	fnSql func(context.Context, *Queries, int64, *UploadImageTransactionResult) error,
	fnMongo func(mongo.SessionContext, string, string) error,
	) error {
	// postgres
	sqlTransaction, err := store.UserDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sqlQuerie := New(sqlTransaction)
	err = fnSql(sqlQuerie, arg.AccountID)
	if err != nil {
		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}

	// mongodb
	session, err := store.ImageDBClient.StartSession()
	if err != nil {
		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}
	defer session.EndSession()

	err = session.StartTransaction()
	if err != nil {
		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}

	// TODO: move database into fnMongo, implement uploading image in fnMongo
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		return fnMongo(sessionContext, arg) // Pass file path to fnMongo
	})
	if er != nil {
		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}

	// commit
	err = sqlTransaction.Commit()
	if err != nil {
		return err
	}
	return session.CommitTransaction(ctx)
}
