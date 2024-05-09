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

func callback(sessionContext mongo.SessionContext, fn func(mongo.SessionContext) error) (interface{}, error) {
	err := fn(sessionContext)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (store *Store) execTransaction(ctx context.Context,
	filepath string,
	fnSql func(*Queries) error,
	fnMongo func(mongo.SessionContext, string) error) error {
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
		return fnMongo(sessionContext, filePath) // Pass file path to fnMongo
	})
	if er != nil {
		if rollbackErr := sqlTransaction.Rollback(), rollbackErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, rollbackErr)
		}
		return err
	}
	mongoDB, err := ImageDBClient.Database()

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
