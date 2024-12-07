package spotify

import (
    "log"
    "github.com/caarlos0/env/v6"
    "github.com/joho/godotenv"
)

type Settings struct{
    SpotifyAccountsUrl   string  `env:"SPOTIFY_ACCOUNTS_URL"`
    SpotifyApiUrl        string  `env:"SPOTIFY_API_URL"`
    SpotifyClientId      string  `env:"SPOTIFY_CLIENT_ID"`
    SpotifyClientToken   string  `env:"SPOTIFY_CLIENT_TOKEN"`
    SpotifyUserId        string  `env:"SPOTIFY_USER_ID"`
    Secret               string  `env:SECRET`
}

func Env() Settings{
    godotenv.Load()
    settings := Settings{}
    if err := env.Parse(&settings);err != nil{
        log.Fatal(err)
    }
    return settings
}

