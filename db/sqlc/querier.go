// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	UpdateEmail(ctx context.Context, arg UpdateEmailParams) (Account, error)
	UpdateImageCount(ctx context.Context, arg UpdateImageCountParams) (Account, error)
	UpdateSubscribed(ctx context.Context, arg UpdateSubscribedParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
