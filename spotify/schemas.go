package spotify

type ApiTokenResponse struct{
    AccessToken   string  `json:"access_token"`
    TokenType     string  `json:"token_type"`
    ExpiresIn     int     `json:"expires_in"`
}

type UserPlaylistsResponse struct{
    Items []struct{
        ID      string   `json:"id"`
        Name    string   `json:"name"`
        URI     string   `json:"uri"`
        Tracks  struct {
            Href   string   `json:"href"`
            Total  int      `json:"total"`
        } `json:"tracks"`
    } `json:"items"`
}

type UserPlaylistTracks struct{
    Items []struct{
        Track struct{
            Name string
            Album struct{
                Name string
                Artists []struct{
                    Name string
                } `json:"artists"`
            } `json:"album"`
        } `json:"track"`
    } `json:"items"`
}
