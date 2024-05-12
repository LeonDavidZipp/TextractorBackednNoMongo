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
	mongodb.MongoOperator
	UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) UploadImageTransactionResult
}

type SQLMongoStore struct {
	*db.Queries
	*mongodb.MongoOperations
	UserDB       *sql.DB
	ImageDB      *mongo.Database
	MongoSession *mongo.Session
}

func NewStore(userDB *sql.DB, session *mongo.Session, imageDB *mongo.Database) Store {
	return &SQLMongoStore{
		Queries: db.New(userDB),
		MongoOperations: mongodb.NewMongo(session, imageDB),
		UserDB: userDB,
		ImageDB: imageDB,
	}
}

type UploadImageTransactionResult struct {
	Image    mongodb.Image `json:"image"`
	Uploader db.Account    `json:"uploader"`
}

// Uploads the data to postgres accounts table
func (store *SQLMongoStore) uploadToPostgres(ctx context.Context, accountID int64, result *UploadImageTransactionResult) error {
	var err error
	result.Uploader, err = store.UpdateImageCount(ctx, db.UpdateImageCountParams{
		Amount: 1,
		ID: accountID,
	})
	if err != nil {
		return err
	}
	return nil
}

// Uploads the image to mongodb table
func (store *SQLMongoStore) uploadToMongo(ctx context.Context, arg UploadImageTransactionParams, result *UploadImageTransactionResult) error {
	var err error
	result.Image, err = store.InsertImage(ctx, mongodb.InsertImageParams{
		AccountID: arg.AccountID,
		Text: arg.Text,
		Link: arg.Link,
		Image64: arg.Image64,
	})
	if err != nil {
		return err
	}
	return nil
}

// execTransaction creates a "rollback-able" transaction.
// It takes as arguments
// 	the context
// 	the filepath to the file to upload
// 	the transcribed text
// 	a function uploading the account data to accounts postgresql table
// 	a function uploading image to mongodb table
// func (store *SQLMongoStore) execTransaction(
// 	ctx context.Context,
// 	arg UploadImageTransactionParams,
// 	result *UploadImageTransactionResult,
// 	fnSql func(context.Context, *db.Queries, int64, *UploadImageTransactionResult) error,
// 	fnMongo func(context.Context, string, string, *UploadImageTransactionResult) error,
// ) error {
// 	// postgres
// 	sqlTransaction, err := store.UserDB.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	sqlQuerie := db.New(sqlTransaction)
// 	err = fnSql(ctx, sqlQuerie, arg.AccountID, result)
// 	if err != nil {
// 		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
// 			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
// 		}
// 		return err
// 	}

// 	// mongodb
// 	session, err := store.ImageDB.Client().StartSession()
// 	if err != nil {
// 		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
// 			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
// 		}
// 		return err
// 	}
// 	defer session.EndSession(ctx)

// 	err = session.StartTransaction()
// 	if err != nil {
// 		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
// 			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
// 		}
// 		return err
// 	}

// 	// TODO: move database into fnMongo, implement uploading image in fnMongo
// 	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
// 		return fnMongo(ctx, arg, sessionContext, result)
// 	})
// 	if err != nil {
// 		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
// 			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
// 		}
// 		return err
// 	}

// 	// commit
// 	err = sqlTransaction.Commit()
// 	if err != nil {
// 		return err
// 	}
// 	return session.CommitTransaction(ctx)
// }

func (store *SQLMongoStore) execTransaction(
	ctx context.Context,
	arg UploadImageTransactionParams,
	result *UploadImageTransactionResult,
	fnSql func(context.Context, *db.Queries) error,
	fnMongo func(context.Context, arg ) error,
) error {
	// postgres
	sqlTransaction, err := store.UserDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sqlQuerie := db.New(sqlTransaction)
	err = fnSql(ctx, sqlQuerie)
	if err != nil {
		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}

	// mongodb
	session, err := store.ImageDB.Client().StartSession()
	if err != nil {
		if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}
	defer session.EndSession(ctx)

	// err = session.StartTransaction()
	// if err != nil {
	// 	if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
	// 		return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
	// 	}
	// 	return err
	// }

	// err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
	// 	return fnMongo(ctx, arg)
	// })
	// if err != nil {
	// 	if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
	// 		return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
	// 	}
	// 	return err
	// }

	err = mongo.WithSession(ctx, session, func(seCtx mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
				return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
			}
			return err
		}
		err = fnMongo(ctx, arg)
	})
	if err != nil {
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

type UploadImageTransactionParams struct {
	AccountID int64  `json:"account_id"`
	Image64   string `json:"image_64"`
	Text      string `json:"text"`
	Link      string `json:"link"`
}

// Upload Image handles uploading the necessary data and image to the databases.
func (store *SQLMongoStore) UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) UploadImageTransactionResult {
	var result UploadImageTransactionResult
	
	err := store.execTransaction(
		ctx,
		arg,
		func(q *db.Queries) error {
			return store.uploadToPostgres(ctx, arg.AccountID, &result)
			},
		func(q *db.Queries) error {
			return store.uploadToMongo(ctx, arg, &result)
		},
	)
	if err != nil {
		return UploadImageTransactionResult{}
	}
}
