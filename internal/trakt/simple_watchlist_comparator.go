package trakt

import (
	"context"
)

type WatchlistComparison struct {
	// ItemsInWatchlist1NotInWatchlist2 are the items that are in the watchlist of user1 but not in the watchlist of user2
	ItemsInWatchlist1NotInWatchlist2 Movies
	// ItemsInWatchlist2NotInWatchlist1 are the items that are in the watchlist of user2 but not in the watchlist of user1
	ItemsInWatchlist2NotInWatchlist1 Movies
	// ItemsInBothWatchlists are the items that are in both watchlists
	ItemsInBothWatchlists Movies
	// Similarity is the similarity between the two watchlists
	Similarity float64
}

type SimpleWatchlistComparator struct{}

func NewSimpleWatchlistComparator() *SimpleWatchlistComparator {
	return &SimpleWatchlistComparator{}
}

// CompareWatchlists compares two watchlists and returns the comparison.
func (SimpleWatchlistComparator) CompareWatchlists(
	_ context.Context,
	user1, user2 Movies,
) (WatchlistComparison, error) {
	var itemsInWatchlist1NotInWatchlist2 Movies
	var itemsInWatchlist2NotInWatchlist1 Movies
	var itemsInBothWatchlists Movies

	// Find items that are in user1's watchlist but not in user2's watchlist
	for _, item := range user1 {
		if !contains(user2, item) {
			itemsInWatchlist1NotInWatchlist2 = append(itemsInWatchlist1NotInWatchlist2, item)
		} else {
			itemsInBothWatchlists = append(itemsInBothWatchlists, item)
		}
	}

	// Find items that are in user2's watchlist but not in user1's watchlist
	for _, item := range user2 {
		if !contains(user1, item) {
			itemsInWatchlist2NotInWatchlist1 = append(itemsInWatchlist2NotInWatchlist1, item)
		}
	}

	// Jaccard similarity between the two watchlists
	similarity := float64(len(itemsInBothWatchlists)) / float64(len(user1)+len(user2)-len(itemsInBothWatchlists))

	return WatchlistComparison{
		ItemsInWatchlist1NotInWatchlist2: itemsInWatchlist1NotInWatchlist2,
		ItemsInWatchlist2NotInWatchlist1: itemsInWatchlist2NotInWatchlist1,
		ItemsInBothWatchlists:            itemsInBothWatchlists,
		Similarity:                       similarity,
	}, nil
}

func contains(movies Movies, movie *Movie) bool {
	for _, itMovie := range movies {
		if itMovie.ID == movie.ID {
			return true
		}
	}

	return false
}
