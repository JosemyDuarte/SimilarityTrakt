package main

import (
	"net/http"

	"golang.org/x/oauth2"

	"MovieTinder/internal/trakt"
)

func main() {
	settings := BuildSettings()

	conf := &oauth2.Config{
		ClientID:     settings.ClientID,
		ClientSecret: settings.ClientSecret,
		Scopes:       []string{"public"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  settings.Domain + "/oauth/authorize",
			TokenURL: settings.Domain + "/oauth/token",
		},
	}

	tokenStorage := trakt.NewInMemoryUserAccessTokenStorage()
	traktService := trakt.NewService(
		tokenStorage,
		trakt.NewInMemoryWatchlistService(settings),
		trakt.NewSimpleWatchlistComparator(),
	)

	traktTokenVerifier := trakt.NewTokenVerifierMiddleware(conf, settings, tokenStorage)

	server := trakt.NewHTTPServer(traktService)

	// Handle watchlist requests
	http.HandleFunc(
		"/trakt/watchlist",
		trakt.EnableCorsMiddleware(traktTokenVerifier.Verify(server.GetWatchlist())),
	)

	// Handle similarity requests
	http.HandleFunc(
		"/trakt/similarity",
		trakt.EnableCorsMiddleware(traktTokenVerifier.Verify(server.CalculateSimilarity(traktService))),
	)

	println("Starting web traktService on port 8080")
	http.ListenAndServe("localhost:8080", nil)
}
