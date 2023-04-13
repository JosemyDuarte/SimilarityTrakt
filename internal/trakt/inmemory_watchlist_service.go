package trakt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type InMemoryWatchlistService struct {
	conf     *oauth2.Config
	settings *Settings
}

func NewInMemoryWatchlistService(settings *Settings) *InMemoryWatchlistService {
	conf := &oauth2.Config{
		ClientID:     settings.ClientID,
		ClientSecret: settings.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", settings.Domain),
			TokenURL: fmt.Sprintf("%s/oauth/token", settings.Domain),
		},
	}

	return &InMemoryWatchlistService{
		conf:     conf,
		settings: settings,
	}
}

func (s InMemoryWatchlistService) GetWatchlist(ctx context.Context, request *GetWatchlistRequest) (Movies, error) {
	if request.Username == "" {
		return nil, fmt.Errorf("username is missing")
	}

	if request.AccessToken == "" {
		return nil, fmt.Errorf("access token is missing")
	}

	token := &oauth2.Token{AccessToken: request.AccessToken}
	httpClient := s.conf.Client(ctx, token)

	watchlistURL := fmt.Sprintf("%s/users/%s/watchlist/movies", s.settings.Domain, request.Username)

	req, err := http.NewRequest("GET", watchlistURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trakt-api-version", s.settings.APIVersion)
	req.Header.Set("trakt-api-key", s.settings.ClientID)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %w", err)
	}
	defer res.Body.Close()

	var watchlist Movies

	err = json.NewDecoder(res.Body).Decode(&watchlist)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return watchlist, nil
}
