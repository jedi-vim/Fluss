package spotify

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func GetPlaylists(accessToken string) UserPlaylistsResponse{
    env := Env()

    request, _ := http.NewRequest(http.MethodGet, 
                                  env.SpotifyApiUrl+"/users/"+env.SpotifyUserId+"/playlists",
                                  nil)
    request.Header.Add("Authorization", "Bearer " + accessToken)
    request.Header.Add("Content-Type", "application/json")

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        log.Fatal(err)
    }
    bodyResponse, _ := io.ReadAll(response.Body)
    if response.StatusCode != 200{
        log.Fatalf("Erro na requisicao:\n %s", bodyResponse)
    }
    defer response.Body.Close()

    var responseJson UserPlaylistsResponse
    if jsonErr := json.Unmarshal(bodyResponse, &responseJson); jsonErr != nil{
        log.Fatal(jsonErr)
    }
    return responseJson
}

func GetPlaylistTracks(accessToken string, playlistID string) (trackList UserPlaylistTracks){
    env := Env()
    request, _ := http.NewRequest(http.MethodGet,
                                  env.SpotifyApiUrl+"/playlists/"+playlistID+"/tracks",
                                  nil)
    request.Header.Add("Authorization", "Bearer " + accessToken)
    request.Header.Add("Content-Type", "application/json")

    response, err := http.DefaultClient.Do(request)
    if err!=nil{
        log.Fatal(err)
    }
    defer response.Body.Close()
    
    bodyResponse, _ := io.ReadAll(response.Body)
    if response.StatusCode != 200{
        log.Printf("Erro na requisicao:\n %s", bodyResponse)
        os.Exit(1)
    }

    var jsonResponse UserPlaylistTracks
    if err := json.Unmarshal(bodyResponse, &jsonResponse); err!=nil{
        log.Fatal(err)
    }
    return jsonResponse
}


