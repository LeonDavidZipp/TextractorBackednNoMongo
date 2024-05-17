package db

import (
	"context"
	"fmt"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
	"go.mongodb.org/mongo-driver/mongo"
)


func (store *DBStore) execTransaction(
	ctx context.Context,
	fnS3 func(*bucket.Client) error,
	rollbackS3 func(*bucket.Client) error,
	fnSql func(*db.Queries) error,
	fnMongo func(*mongodb.MongoOperations) error,
) error {
	// s3: AUTOMATICALLY COMMITED WHEN SUCCESSFUL

	if err := fnS3(store.Client); err != nil {
		return err
	}
	
	// postgres
	sqlTransaction, err := store.UserDB.BeginTx(ctx, nil)
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

	// mongodb
	session, err := store.ImageDB.Client().StartSession()
	if err != nil {
		if sqlRbErr := sqlTransaction.Rollback(); sqlRbErr != nil {
			if s3RbErr := rollbackS3(store.Client); s3RbErr != nil {
				return fmt.Errorf("transaction err: %v; sql rollback err: %v; s3 rollback err: %v", err, sqlRbErr, s3RbErr)
			}
			return fmt.Errorf("transaction err: %v; sql rollback err: %v", err, sqlRbErr)
		}
		if s3RbErr := rollbackS3(store.Client); s3RbErr != nil {
			return fmt.Errorf("transaction err: %v; s3 rollback err: %v", err, s3RbErr)
		}
		return err
	}

	err = mongo.WithSession(ctx, session, func(seCtx mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			if sqlRbErr := sqlTransaction.Rollback(); sqlRbErr != nil {
				if s3RbErr := rollbackS3(store.Client); s3RbErr != nil {
					return fmt.Errorf("transaction err: %v; sql rollback err: %v; s3 rollback err: %v", err, sqlRbErr, s3RbErr)
				}
				return fmt.Errorf("transaction err: %v; sql rollback err: %v", err, sqlRbErr)
			}
			if s3RbErr := rollbackS3(store.Client); s3RbErr != nil {
				return fmt.Errorf("transaction err: %v; s3 rollback err: %v", err, s3RbErr)
			}
			return err
		}

		err = fnMongo(store.MongoOperations)
		if err != nil {
			s3RbErr := rollbackS3(store.Client)
			sqlRbErr := sqlTransaction.Rollback()
			mongoRbErr := session.AbortTransaction(ctx)
			if sqlRbErr != nil || s3RbErr != nil || mongoRbErr != nil {
				if sqlRbErr != nil && mongoRbErr != nil && s3RbErr != nil {
					return fmt.Errorf("mongo transaction err: %v; sql rollback err: %v; mongo rollback err: %v; s3 rollback err: %v", err, sqlRbErr, mongoRbErr, s3RbErr)
				} else if sqlRbErr != nil && s3RbErr != nil {
					return fmt.Errorf("mongo transaction err: %v; sql rollback err: %v; s3 rollback err: %v", err, sqlRbErr, s3RbErr)
				} else if sqlRbErr != nil && mongoRbErr != nil {
					return fmt.Errorf("mongo transaction err: %v; sql rollback err: %v; mongo rollback err: %v", err, sqlRbErr, mongoRbErr)
				} else if mongoRbErr != nil && s3RbErr != nil {
					return fmt.Errorf("mongo transaction err: %v; mongo rollback err: %v; s3 rollback err: %v", err, mongoRbErr, s3RbErr)
				} else if sqlRbErr != nil {
					return fmt.Errorf("mongo transaction err: %v; sql rollback err: %v", err, sqlRbErr)
				} else if mongoRbErr != nil {
					return fmt.Errorf("mongo transaction err: %v; mongo rollback err: %v", err, mongoRbErr)
				} else if s3RbErr != nil {
					return fmt.Errorf("mongo transaction err: %v; s3 rollback err: %v", err, s3RbErr)
				}
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
		if s3RbErr := rollbackS3(store.Client); s3RbErr != nil {
			if mongoRbErr := session.AbortTransaction(ctx); mongoRbErr != nil {
				return fmt.Errorf("sql transaction err: %v; s3 rollback err: %v; mongo rollback err: %v", err, s3RbErr, mongoRbErr)
			}
			return fmt.Errorf("sql transaction err: %v; s3 rollback err: %v", err, s3RbErr)
		}
		return err
	}
	err = session.CommitTransaction(ctx)
	if err != nil {
		if s3RbErr := rollbackS3(store.Client); s3RbErr != nil {
			return fmt.Errorf("mongo transaction err: %v; s3 rollback err: %v", err, s3RbErr)
		}
		return fmt.Errorf("mongo transaction err: %v", err)
	}
	return nil
}
