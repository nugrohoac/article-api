package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestSuite .
type TestSuite struct {
	suite.Suite
	Client *redis.Client
}

// SetupSuite .
func (t *TestSuite) SetupSuite() {
	t.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

// TearDownTest .
func (t *TestSuite) TearDownTest() {
	err := t.Client.FlushDB(context.Background()).Err()
	require.NoError(t.T(), err)
}
