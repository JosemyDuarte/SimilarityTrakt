package trakt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// EnableCorsMiddleware enables CORS for the given handler.
func EnableCorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Update this to be more secure
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		next(w, r)
	}
}

// TokenVerifierMiddleware verifies the oauth2 token and adds the username to the context.
type TokenVerifierMiddleware struct {
	conf       *oauth2.Config
	settings   *Settings
	tokenCache UserTokenStorage
}

func NewTokenVerifierMiddleware(
	conf *oauth2.Config,
	settings *Settings,
	tokenCache UserTokenStorage,
) *TokenVerifierMiddleware {
	return &TokenVerifierMiddleware{conf: conf, settings: settings, tokenCache: tokenCache}
}

// Verify verifies the oauth2 token and adds the username to the context
// by calling Trakt API.
func (m *TokenVerifierMiddleware) Verify(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "access token is missing", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		username, err := m.getUsername(ctx, accessToken)
		if err != nil {
			http.Error(w, "failed to get username", http.StatusUnauthorized)
			return
		}

		err = m.tokenCache.StoreAccessToken(ctx, username, accessToken)
		if err != nil {
			http.Error(w, "failed to store access token", http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, "username", username)
		ctx = context.WithValue(ctx, "token", accessToken)
		next(w, r.WithContext(ctx))
	}
}

// getUsername gets the username from the Trakt API using the access token.
func (m *TokenVerifierMiddleware) getUsername(ctx context.Context, accessToken string) (string, error) {
	token := &oauth2.Token{AccessToken: accessToken}
	httpClient := m.conf.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api-staging.trakt.tv/users/me", nil)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trakt-api-version", m.settings.APIVersion)
	req.Header.Set("trakt-api-key", m.settings.ClientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get username, status code: %d", resp.StatusCode)
	}

	var user struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	return user.Username, nil
}
