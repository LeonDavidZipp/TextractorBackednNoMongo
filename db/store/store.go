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
	UserDB  *sql.DB
	ImageDB *mongo.Database
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

	err = mongo.WithSession(ctx, session, func(seCtx mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			if rollbackErr := sqlTransaction.Rollback(); rollbackErr != nil {
				return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
			}
			return err
		}
		err = fnMongo(store.MongoOperations)
		if err != nil {
			sqlRollbackErr := sqlTransaction.Rollback()
			mongoRollbackErr := session.AbortTransaction(ctx)
			if sqlRollbackErr != nil {
				if mongoRollbackErr != nil {
					return fmt.Errorf("mongo transaction err: %v; sql rollback err: %v; mongo rollback err: %v", err, sqlRollbackErr, mongoRollbackErr)
				}
				return fmt.Errorf("mongo transaction err: %v; sql rollback err: %v", err, sqlRollbackErr)
			} else if mongoRollbackErr != nil {
				return fmt.Errorf("mongo transaction err: %v; mongo rollback err: %v", err, mongoRollbackErr)
			} else {
				return err
			}
		}
		return nil
		},
	)

	// commit
	err = sqlTransaction.Commit()
	if err != nil {
		if rollbackErr := session.AbortTransaction(ctx); rollbackErr != nil {
			return fmt.Errorf("mongo transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}
	err = session.CommitTransaction(ctx)
	if err != nil {
		return fmt.Errorf("mongo transaction err: %v", err)
	}
	return nil
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
	var uploader db.Account
	var image mongodb.Image

	err := store.execTransaction(
		ctx,
		func(q *db.Queries) error {
			var err error
			uploader, err = q.UpdateImageCount(ctx, db.UpdateImageCountParams{
				Amount: 1,
				ID: arg.AccountID,
			})
			return err
		},
		func(op *mongodb.MongoOperations) error {
			var err error
			image, err = op.InsertImage(ctx, mongodb.InsertImageParams{
				AccountID: arg.AccountID,
				Text: arg.Text,
				Link: arg.Link,
				Image64: arg.Image64,
			})
			return err
		},
	)
	if err != nil {
		return UploadImageTransactionResult{}, err
	}

	result := UploadImageTransactionResult{
		Image: image,
		Uploader: uploader,
	}

	return result, nil
}
