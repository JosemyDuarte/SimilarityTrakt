package trakt

type Username string

type User struct {
	Username    Username `json:"username"`
	Watchlist   *Movies  `json:"watchlist"`
	AccessToken string   `json:"access_token"`
}
