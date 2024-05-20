package db

import (
	"context"
	"testing"
	"time"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/LeonDavidZipp/Textractor/util"
)


func createRandomUser(t *testing.T) User {
	name := util.RandomString(8)

	ctx := context.Background()
	user, err := testUserQueries.CreateUser(ctx, name)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	ctx := context.Background()
	user2, err := testUserQueries.GetUser(ctx, user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.ImageCount, user2.ImageCount)
	require.Equal(t, user1.Subscribed, user2.Subscribed)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, 10 * time.Second)
}

func TestUpdateSubscribed(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateSubscribedParams{
		ID : user1.ID,
		Subscribed : true,
	}

	ctx := context.Background()
	user2, err := testUserQueries.UpdateSubscribed(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.ImageCount, user2.ImageCount)
	require.Equal(t, arg.Subscribed, user2.Subscribed)
}

func TestUpdateImageCount(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateImageCountParams{
		ID : user1.ID,
		Amount : 10,
	}

	ctx := context.Background()
	user2, err := testUserQueries.UpdateImageCount(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.ImageCount + arg.Amount, user2.ImageCount)
	require.Equal(t, user1.Subscribed, user2.Subscribed)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	ctx := context.Background()
	err := testUserQueries.DeleteUser(ctx, user1.ID)
	require.NoError(t, err)

	user2, err := testUserQueries.GetUser(ctx, user1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i :=0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit : 5,
		Offset : 5,
	}

	ctx := context.Background()
	users, err := testUserQueries.ListUsers(ctx, arg)

	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
