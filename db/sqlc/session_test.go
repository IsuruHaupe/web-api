package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createFakeSession(t *testing.T) (Session, uuid.UUID) {
	user := CreateRandomUser(t)
	ID := uuid.New()
	arg := CreateSessionParams{
		ID:           ID,
		Username:     user.Username,
		SessionToken: "Token",
		UserAgent:    "Agent",
		ClientIp:     "IP",
		IsBlocked:    false,
		ExpiresAt:    time.Now(),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	return session, ID
}
func TestCreateSession(t *testing.T) {
	_, ID := createFakeSession(t)

	err := testQueries.DeleteSession(context.Background(), ID)
	require.NoError(t, err)
}

func TestGetSession(t *testing.T) {
	session1, ID := createFakeSession(t)
	session2, err := testQueries.GetSession(context.Background(), ID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.Username, session2.Username)
	require.Equal(t, session1.SessionToken, session2.SessionToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)
	require.WithinDuration(t, session1.ExpiresAt, session2.ExpiresAt, time.Second)

	err = testQueries.DeleteSession(context.Background(), ID)
	require.NoError(t, err)
}

func TestDeleteSession(t *testing.T) {
	_, ID := createFakeSession(t)
	err := testQueries.DeleteSession(context.Background(), ID)
	require.NoError(t, err)

	session2, err := testQueries.GetSession(context.Background(), ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, session2)
}
