package api

import (
	"os"
	"testing"
	"time"

	"github.com/IsuruHaupe/web-api/auth"
	"github.com/IsuruHaupe/web-api/config"
	"github.com/IsuruHaupe/web-api/db/postgres"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, database postgres.Database) *Server {
	config := config.Config{
		TokenSymmetricKey:   auth.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, database)
	require.NoError(t, err)

	return server
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
