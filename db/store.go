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
	UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error)
}

type SQLMongoStore struct {
	*db.Queries
	*mongodb.MongoOperations
	UserDB       *sql.DB
	ImageDB      *mongo.Database
}

func NewStore(
	userDB *sql.DB,
	imageDB *mongo.Database) Store {
	return &SQLMongoStore{
		Queries: db.New(userDB),
		MongoOperations: mongodb.NewMongo(imageDB),
		UserDB: userDB,
		ImageDB: imageDB,
	}
}

func (store *SQLMongoStore) execTransaction(
	ctx context.Context,
	fnSql func(*db.Queries) error,
	fnMongo func(*mongodb.MongoOperations) error,
) error {
	// postgres
	sqlTransaction, err := store.UserDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sqlQuerie := db.New(sqlTransaction)
	err = fnSql(sqlQuerie)
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

	err = mongo.WithSession(ctx, session, func(seCtx mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
				return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
			}
			return err
		}
		if err != nil {
			if rollbackErr := session.AbortTransaction(ctx); rollbackErr != nil {
				return fmt.Errorf("mongo transaction err: %v; rollback err: %v", err, rollbackErr)
			}
			return err
		}
		return nil
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
	Text      string `json:"text"`
	Link      string `json:"link"`
	Image64   string `json:"image_64"`
}

type UploadImageTransactionResult struct {
	Image    mongodb.Image `json:"image"`
	Uploader db.Account    `json:"uploader"`
}

// Upload Image handles uploading the necessary data and image to the databases.
func (store *SQLMongoStore) UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error) {
	var result UploadImageTransactionResult
	
	err := store.execTransaction(
		ctx,
		func(q *db.Queries) error {
			var err error
			result.Uploader, err = q.UpdateImageCount(ctx, db.UpdateImageCountParams{
				Amount: 1,
				ID: arg.AccountID,
			})
			if err != nil {
				return err
			}
			return nil
		},
		func(op *mongodb.MongoOperations) error {
			var err error
			result.Image, err = op.InsertImage(ctx, mongodb.InsertImageParams{
				AccountID: arg.AccountID,
				Text: arg.Text,
				Link: arg.Link,
				Image64: arg.Image64,
			})
			if err != nil {
				return err
			}
			return nil
		},
	)
	return result, err
}
