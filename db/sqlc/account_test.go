package db

import (
	"context"
	"testing"
	"time"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/LeonDavidZipp/Textractor/util"
)


func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner : util.RandomName(),
		Email : util.RandomEmail(),
		GoogleID : sql.NullString{},
		FacebookID : sql.NullString{},
	}

	ctx := context.Background()
	account, err := testAccountQueries.CreateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Email, account.Email)
	require.Equal(t, arg.GoogleID, account.GoogleID)
	require.Equal(t, arg.FacebookID, account.FacebookID)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	ctx := context.Background()
	account2, err := testAccountQueries.GetAccount(ctx, account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.GoogleID, account2.GoogleID)
	require.Equal(t, account1.FacebookID, account2.FacebookID)
	require.Equal(t, account1.ImageCount, account2.ImageCount)
	require.Equal(t, account1.Subscribed, account2.Subscribed)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 10 * time.Second)
}

func TestUpdateEmail(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateEmailParams{
		ID : account1.ID,
		Email : "leon@example.com",
	}

	ctx := context.Background()
	account2, err := testAccountQueries.UpdateEmail(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.GoogleID, account2.GoogleID)
	require.Equal(t, account1.FacebookID, account2.FacebookID)
	require.Equal(t, account1.ImageCount, account2.ImageCount)
	require.Equal(t, account1.Subscribed, account2.Subscribed)
	require.Equal(t, arg.Email, account2.Email)
}

func TestUpdateSubscribed(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateSubscribedParams{
		ID : account1.ID,
		Subscribed : true,
	}

	ctx := context.Background()
	account2, err := testAccountQueries.UpdateSubscribed(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.GoogleID, account2.GoogleID)
	require.Equal(t, account1.FacebookID, account2.FacebookID)
	require.Equal(t, account1.ImageCount, account2.ImageCount)
	require.Equal(t, arg.Subscribed, account2.Subscribed)
}

func TestUpdateImageCount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateImageCountParams{
		ID : account1.ID,
		Amount : 10,
	}

	ctx := context.Background()
	account2, err := testAccountQueries.UpdateImageCount(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.GoogleID, account2.GoogleID)
	require.Equal(t, account1.FacebookID, account2.FacebookID)
	require.Equal(t, account1.ImageCount + arg.Amount, account2.ImageCount)
	require.Equal(t, account1.Subscribed, account2.Subscribed)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	ctx := context.Background()
	err := testAccountQueries.DeleteAccount(ctx, account1.ID)
	require.NoError(t, err)

	account2, err := testAccountQueries.GetAccount(ctx, account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i :=0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit : 5,
		Offset : 5,
	}

	ctx := context.Background()
	accounts, err := testAccountQueries.ListAccounts(ctx, arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
