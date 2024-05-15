package db

import (
	"context"
	"fmt"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"go.mongodb.org/mongo-driver/mongo"
)


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