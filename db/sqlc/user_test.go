package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/IsuruHaupe/web-api/auth"
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := auth.HashPassword("password")
	require.NoError(t, err)
	fullname := randomdata.FullName(randomdata.Female)
	args := CreateUserParams{
		Username:       fullname,
		HashedPassword: hashedPassword,
		Fullname:       fullname,
		Email:          randomdata.Email(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.Fullname, user.Fullname)
	require.Equal(t, args.Email, user.Email)

	require.True(t, user.PasswordLastChanged.IsZero())
	require.NotZero(t, user.CreateAt)

	return user
}
func TestCreateUser(t *testing.T) {
	user := CreateRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordLastChanged, user2.PasswordLastChanged, time.Second)
	require.WithinDuration(t, user1.CreateAt, user2.CreateAt, time.Second)

	err = testQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}
