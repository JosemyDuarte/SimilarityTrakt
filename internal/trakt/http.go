package trakt

import (
	"encoding/json"
	"net/http"
)

type HTTPServer struct {
	TraktService *Service
}

func NewHTTPServer(traktService *Service) *HTTPServer {
	return &HTTPServer{TraktService: traktService}
}

// GetWatchlist returns the watchlist of a user.
func (s HTTPServer) GetWatchlist() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value("username").(string)
		token := r.Context().Value("token").(string)

		watchlist, err := s.TraktService.GetWatchlist(
			r.Context(),
			&GetWatchlistRequest{AccessToken: token, Username: username},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(watchlist)
	}
}

// CalculateSimilarity calculates the similarity between two users.
func (s HTTPServer) CalculateSimilarity(traktService *Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value("username").(string)
		accessToken := r.Context().Value("token").(string)

		var matchRequest SimilarityRequest

		err := json.NewDecoder(r.Body).Decode(&matchRequest)
		if err != nil || matchRequest.OtherUsername == "" {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		matchRequest.AccessToken = accessToken
		matchRequest.Username = username

		matches, err := traktService.Similarity(r.Context(), &matchRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matches)
	}
}
