package main

import (
    "fmt"
    "log"

    "github.com/alecthomas/kong"

    "spotify-crawler/google"
    "spotify-crawler/spotify"
)

type PlaylistsCmd struct{
}

func (r * PlaylistsCmd) Run(ctx *kong.Context)error{
    accessToken := spotify.GetAccessToken()
    spotify.GetPlaylists(accessToken)
    return nil
}

type MigrateCmd struct{
    PlaylistID string `arg required help:"Spotify PlaylistID"`
    YoutubePlaylistName string `arg required help:"Youtube playlist name"`
}

func (m *MigrateCmd) Run(ctx *kong.Context)error{ 
    playlistInfo := spotify.GetPlaylistTracks(spotify.GetAccessToken(), m.PlaylistID)
    var trackList []string
    for _, item := range playlistInfo.Items{
        track := item.Track 
        trackFullName := fmt.Sprintf("%s %s", track.Name, track.Album.Artists[0].Name)
        log.Println(trackFullName)
        trackList = append(trackList, trackFullName)
    }
    googleEnv := google.Env()
    playlist := google.CreatePlaylist(m.YoutubePlaylistName, &googleEnv)
    google.PopulatePlaylist(playlist, trackList, &googleEnv)
    return nil
}

type CLI struct{
    Playlists PlaylistsCmd `cmd:"" help:"Listar Playlists"`
    Migrate   MigrateCmd   `cmd:"" help:"Migrar Playlist para o Youtube"`
}
