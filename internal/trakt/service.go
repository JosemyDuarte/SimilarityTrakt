package trakt

import (
	"context"
	"fmt"
)

// UserTokenStorage is the interface for storing user access tokens.
//
//go:generate mockgen -source=./service.go -destination=./mocks/service.go -package=mocks UserTokenStorage
type UserTokenStorage interface {
	StoreAccessToken(ctx context.Context, userID string, token string) error
	GetAccessToken(ctx context.Context, username string) (string, error)
}

// WatchlistService is the interface for the watchlist service.
//
//go:generate mockgen -source=./service.go -destination=./mocks/service.go -package=mocks WatchlistService
type WatchlistService interface {
	GetWatchlist(ctx context.Context, request *GetWatchlistRequest) (Movies, error)
}

// WatchlistComparator is the interface for comparing watchlists of two users.
//
//go:generate mockgen -source=./service.go -destination=./mocks/service.go -package=mocks WatchlistComparator
type WatchlistComparator interface {
	CompareWatchlists(ctx context.Context, user1, user2 Movies) (WatchlistComparison, error)
}

type Service struct {
	userTokenStorage    UserTokenStorage
	watchlistService    WatchlistService
	watchlistComparator WatchlistComparator
}

func NewService(
	userTokenStorage UserTokenStorage,
	watchlistService WatchlistService,
	watchlistComparator WatchlistComparator,
) *Service {
	return &Service{
		userTokenStorage:    userTokenStorage,
		watchlistService:    watchlistService,
		watchlistComparator: watchlistComparator,
	}
}

type AuthRequest struct {
	Username string `json:"username"`
}

type AuthCallback struct {
	// State is the state that was sent to the authorization endpoint
	State string `json:"state"`
	// Code is the authorization code that can be exchanged for an access token
	Code string `json:"code,omitempty"`
}

type AuthResponse struct {
	// URL is the URL to redirect the user to
	URL string `json:"url"`
}

type GetWatchlistRequest struct {
	// Username is the username of the user to get the watchlist for
	Username string `json:"username"`
	// AccessToken is the access token for the user
	AccessToken string `json:"access_token"`
}

type GetWatchlistResponse struct {
	// Items is the list of items in the watchlist
	Items Movies `json:"movies"`
}

// GetWatchlist returns the watchlist for the user with the given access token.
func (s *Service) GetWatchlist(ctx context.Context, req *GetWatchlistRequest) (*GetWatchlistResponse, error) {
	watchlist, err := s.watchlistService.GetWatchlist(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", err)
	}

	return &GetWatchlistResponse{Items: watchlist}, nil
}

// SimilarityRequest is the request to check how related two users are.
type SimilarityRequest struct {
	// AccessToken is the access token for the user
	AccessToken string `json:"access_token"`
	// Username is the username of the user to check the similarity for
	Username string `json:"username"`
	// OtherUsername is the username of the other user to compare with
	OtherUsername string `json:"other_username"`
}

// MatchResponse shows how related two users are.
type MatchResponse struct {
	// Score is the score of how related the two users are
	Score float64 `json:"score"`
}

// Similarity checks how similar two users are based on their watchlists.
func (s *Service) Similarity(ctx context.Context, request *SimilarityRequest) (*MatchResponse, error) {
	// Fetch access token for the other user
	otherUserAccessToken, err := s.userTokenStorage.GetAccessToken(ctx, request.OtherUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to find user %s: %w", request.OtherUsername, err)
	}

	// Fetch watchlist for the other user
	otherUserWatchlist, err := s.watchlistService.GetWatchlist(ctx, &GetWatchlistRequest{
		AccessToken: otherUserAccessToken,
		Username:    request.OtherUsername,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", err)
	}

	// Fetch watchlist for user
	watchlist, err := s.watchlistService.GetWatchlist(ctx, &GetWatchlistRequest{
		AccessToken: request.AccessToken,
		Username:    request.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", err)
	}

	// Compare watchlists
	comparison, err := s.watchlistComparator.CompareWatchlists(ctx, watchlist, otherUserWatchlist)
	if err != nil {
		return nil, fmt.Errorf("failed to compare watchlists: %w", err)
	}

	return &MatchResponse{Score: comparison.Similarity}, nil
}
