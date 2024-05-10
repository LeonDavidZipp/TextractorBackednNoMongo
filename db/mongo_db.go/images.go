package db
import (
	"context"
	"database/sql"
	"fmt"
	mongodb "go.mongodb.org/mongo-driver/mongo"
)

func (store *Store) CreateImage(
	ctx context.Context,
	arg UploadImageTransactionParams
) UploadImageTransactionResult {
	
}

func (store *Store) GetImage(ctx context.Context, id string) (*Image, error) {
	
}
