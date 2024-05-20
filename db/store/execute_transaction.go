package db

import (
	"context"
	"fmt"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
)


func (store *DBStore) execTransaction(
	ctx context.Context,
	fnS3 func(*bucket.Client) error,
	rollbackS3 func(*bucket.Client) error,
	fnSql func(*db.Queries) error,
) error {
	// s3: AUTOMATICALLY COMMITED WHEN SUCCESSFUL
	if err := fnS3(store.Client); err != nil {
		return err
	}
	
	// postgres
	sqlTransaction, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		if s3RbErr := rollbackS3(store.Client); s3RbErr != nil {
			return fmt.Errorf("sql err: %v; rollback err: %v", err, s3RbErr)
		}
		return err
	}

	sqlQuerie := db.New(sqlTransaction)
	err = fnSql(sqlQuerie)
	if err != nil {
		if sqlRbErr := sqlTransaction.Rollback(); sqlRbErr != nil {
			return fmt.Errorf("transaction err: %v; rollback err: %v", err, sqlRbErr)
		}
		return err
	}

	// commit
	err = sqlTransaction.Commit()
	if err != nil {
		if s3RbErr := rollbackS3(store.Client); s3RbErr != nil {
			return fmt.Errorf("sql transaction err: %v; s3 rollback err: %v", err, s3RbErr)
		}
		return err
	}
	return nil
}
