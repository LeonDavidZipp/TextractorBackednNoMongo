package db

import (
	"context"
)


type MongoQuerier interface {

}

var _ MongoQuerier = (*MongoQueries)(nil)
