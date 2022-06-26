package api

import (
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/croixxant/go-sandbox/simplebank/db"
	"github.com/croixxant/go-sandbox/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   gofakeit.Password(true, true, true, false, false, 32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
