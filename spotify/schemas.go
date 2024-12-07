package spotify

type ApiTokenResponse struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
        RefreshToken  string `json:refresh_token`
        Scope         string `json:scope`
}

type UserPlaylistsResponse struct {
	Items []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		URI    string `json:"uri"`
		Tracks struct {
			Href  string `json:"href"`
			Total int    `json:"total"`
		} `json:"tracks"`
	} `json:"items"`
}

type UserPlaylistTracks struct {
	Items []struct {
		Track struct {
			Name  string
			Album struct {
				Name    string
				Artists []struct {
					Name string
				} `json:"artists"`
			} `json:"album"`
		} `json:"track"`
	} `json:"items"`
}

type UserSavedTracksResponse struct {
	Href     string  `json:"href"`
	Limit    int     `json:"limit"`
	Next     *string `json:"next"`
	Offset   int     `json:"offset"`
	Previous *string `json:"previous"`
	Total    int     `json:"total"`
	Items    []struct {
		Track   struct {
			Album struct {
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
			} `json:"album"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Name string `json:"name"`
		} `json:"track"`
	} `json:"items"`
}
