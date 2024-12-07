package spotify

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

func GetLikedSongs(accessToken string) (savedSongs []map[string]string){
    savedSongs = []map[string]string{}
    env := Env()

    continueFecthing := true
    offSet := 0
    for continueFecthing {
        queryParams := url.Values{}
        queryParams.Set("offset", strconv.Itoa(offSet))
        queryParams.Set("limit", strconv.Itoa(50))
        request, _ := http.NewRequest(
            http.MethodGet,
            env.SpotifyApiUrl+"/v1/me/tracks?" + queryParams.Encode(),
            nil)
        request.Header.Add("Authorization", "Bearer " + accessToken)
        request.Header.Add("Content-Type", "application/json")
        response, err := http.DefaultClient.Do(request)
        if err != nil {
            log.Fatal(err)
        }

        bodyResponse, _ := io.ReadAll(response.Body)
        if response.StatusCode != 200{
            log.Fatalf("Erro durante a requisicao:\n %s", bodyResponse)
        }
        defer response.Body.Close()

        var jsonResponse UserSavedTracksResponse
        if err := json.Unmarshal(bodyResponse, &jsonResponse); err != nil{
            log.Fatal(err)
        }

        if len(jsonResponse.Items) == 0 {
            continueFecthing = false
            continue
        }

        for _, track := range jsonResponse.Items {
            trackMap := map[string]string{
                "name": track.Track.Name,
                "artist": track.Track.Artists[0].Name,
                "album": track.Track.Album.Name,
                "released": track.Track.Album.ReleaseDate,
            }
            log.Printf("Artista: %s || Musica: %s || Album: %s", 
                trackMap["artist"], trackMap["name"], trackMap["album"])
            savedSongs = append(savedSongs, trackMap)
        }
        offSet = offSet + 50
    }
    return
}
