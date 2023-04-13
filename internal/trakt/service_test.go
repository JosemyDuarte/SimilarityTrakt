package trakt_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"MovieTinder/internal/trakt"
	"MovieTinder/internal/trakt/mocks"
)

func TestService_GetWatchlist(t *testing.T) {
	ctx := context.Background()
	accessToken := "test_token"
	username := "test_user"
	watchlist := trakt.Movies{
		&trakt.Movie{ID: 1},
		&trakt.Movie{ID: 2},
		&trakt.Movie{ID: 3},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserTokenStorage := mocks.NewMockUserTokenStorage(ctrl)

	mockWatchlistService := mocks.NewMockWatchlistService(ctrl)
	mockWatchlistService.EXPECT().GetWatchlist(
		ctx,
		&trakt.GetWatchlistRequest{Username: username, AccessToken: accessToken},
	).Return(watchlist, nil)

	service := trakt.NewService(
		mockUserTokenStorage,
		mockWatchlistService,
		trakt.NewSimpleWatchlistComparator(),
	)

	req := &trakt.GetWatchlistRequest{Username: username, AccessToken: accessToken}
	resp, err := service.GetWatchlist(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, watchlist, resp.Items)
}

func TestService_CalculateSimilarity(t *testing.T) {
	request := &trakt.SimilarityRequest{
		AccessToken:   "test_token",
		Username:      "test_user",
		OtherUsername: "other_user",
	}
	otherAccessToken := "other_token"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockUserTokenStorage := mocks.NewMockUserTokenStorage(ctrl)
	mockUserTokenStorage.EXPECT().GetAccessToken(ctx, request.OtherUsername).Return(otherAccessToken, nil)

	mockWatchlistService := mocks.NewMockWatchlistService(ctrl)
	mockWatchlistService.EXPECT().GetWatchlist(
		ctx,
		&trakt.GetWatchlistRequest{
			Username:    request.OtherUsername,
			AccessToken: otherAccessToken,
		},
	).Return(trakt.Movies{
		&trakt.Movie{ID: 1},
		&trakt.Movie{ID: 2},
	}, nil)

	mockWatchlistService.EXPECT().GetWatchlist(
		ctx,
		&trakt.GetWatchlistRequest{
			Username:    request.Username,
			AccessToken: request.AccessToken,
		},
	).Return(trakt.Movies{
		&trakt.Movie{ID: 1},
		&trakt.Movie{ID: 2},
	}, nil)

	service := trakt.NewService(
		mockUserTokenStorage,
		mockWatchlistService,
		trakt.NewSimpleWatchlistComparator(),
	)

	resp, err := service.Similarity(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 1.0, resp.Score, "should be 1.0 because both users have the same watchlist")
}
