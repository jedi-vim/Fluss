package main

import (
	"log"
        "spotify-crawler/spotify"
)

func main(){
    accessToken := getAccessTokenFromFile() 
    if accessToken == ""{
        log.Println("Nao Existe token valido em arquivo, vou gerar um novo")
        accessToken = spotify.GenerateNewAccessToken()
        saveToken(accessToken)
    }
    playlistsReferences := spotify.GetPlaylists(accessToken)
    spotify.SavePlaylistsTracks(accessToken, playlistsReferences)
}
