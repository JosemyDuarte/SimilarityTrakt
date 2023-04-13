package trakt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemoryUserAccessTokenStorage_StoreAndGetAccessToken(t *testing.T) {
	storage := NewInMemoryUserAccessTokenStorage()

	// Test StoreAccessToken method
	err := storage.StoreAccessToken(context.Background(), "user1", "token1")
	require.NoError(t, err)

	err = storage.StoreAccessToken(context.Background(), "user2", "token2")
	require.NoError(t, err)

	// Test GetAccessToken method
	token1, err := storage.GetAccessToken(context.Background(), "user1")
	require.NoError(t, err)
	require.Equal(t, "token1", token1)

	token2, err := storage.GetAccessToken(context.Background(), "user2")
	require.NoError(t, err)
	require.Equal(t, "token2", token2)

	// Test non-existing user
	_, err = storage.GetAccessToken(context.Background(), "user3")
	require.Error(t, err)
	require.EqualError(t, err, "token for user user3 not found")
}
