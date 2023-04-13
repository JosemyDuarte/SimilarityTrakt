package trakt

import "time"

type MovieInfo struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	IDs   struct {
		Trakt int    `json:"trakt"`
		Slug  string `json:"slug"`
		IMDb  string `json:"imdb"`
		TMDB  int    `json:"tmdb"`
	} `json:"ids"`
}

type Movie struct {
	Rank      int       `json:"rank"`
	ID        int       `json:"id"`
	ListedAt  time.Time `json:"listed_at"`
	Notes     string    `json:"notes"`
	Type      string    `json:"type"`
	MovieInfo MovieInfo `json:"movie"`
}

type Movies []*Movie
