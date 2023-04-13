package trakt_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"MovieTinder/internal/trakt"
)

func TestSimpleWatchlistComparator_CompareWatchlists(t *testing.T) {
	// Set up test data
	movie1 := &trakt.Movie{ID: 1}
	movie2 := &trakt.Movie{ID: 2}
	movie3 := &trakt.Movie{ID: 3}
	movie4 := &trakt.Movie{ID: 4}

	user1 := trakt.Movies{movie1, movie2, movie3}
	user2 := trakt.Movies{movie2, movie3, movie4}

	// Create comparator and call CompareWatchlists method
	comparator := trakt.NewSimpleWatchlistComparator()
	comparison, err := comparator.CompareWatchlists(context.Background(), user1, user2)

	// Check for errors
	require.NoError(t, err)

	// Check for correct items in each category
	require.Equal(t, trakt.Movies{movie1}, comparison.ItemsInWatchlist1NotInWatchlist2)
	require.Equal(t, trakt.Movies{movie4}, comparison.ItemsInWatchlist2NotInWatchlist1)
	require.Equal(t, trakt.Movies{movie2, movie3}, comparison.ItemsInBothWatchlists)

	// Check for correct similarity calculation
	require.Equal(t, 0.5, comparison.Similarity)
}
