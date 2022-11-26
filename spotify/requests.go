package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func generateBearerToken(spotifyClientId string, spotifyClientToken string) string{
    firstStep := fmt.Sprintf("%s:%s", spotifyClientId, spotifyClientToken)
    firstStepToBytes := []byte(firstStep)
    return base64.StdEncoding.EncodeToString(firstStepToBytes)
}

func GenerateNewAccessToken() string {
    enviroment := Env()
    formData := url.Values{"grant_type": {"client_credentials"}}

    authReq, err := http.NewRequest(
        http.MethodPost, 
        enviroment.SpotifyAccountsUrl+"/api/token", 
        strings.NewReader(formData.Encode()))
    if err != nil{
        log.Fatal(err)
    }

    bearerToken := fmt.Sprintf("Basic %s", generateBearerToken(enviroment.SpotifyClientId, enviroment.SpotifyClientToken))
    authReq.Header.Set( "Content-Type", "application/x-www-form-urlencoded")
    authReq.Header.Add( "Authorization", bearerToken)
    
    log.Println("Fazendo auth-request para o Spotify")
    res, err := http.DefaultClient.Do(authReq)
    if err != nil || res.StatusCode != 200{
        log.Fatal(err)
    }
    defer res.Body.Close()

    bodyResp, _ := io.ReadAll(res.Body)
    log.Printf("Resposta recebida do spotify:\n %s\n", bodyResp)

    var responseJson ApiTokenResponse
    if jsonErr := json.Unmarshal(bodyResp, &responseJson); jsonErr != nil{
        log.Println("Erro no decode da resposta json")
        log.Fatal(jsonErr)
    }
    return responseJson.AccessToken
}

func GetPlaylists(accessToken string) (playlistsUri map[string]string){
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
        log.Printf("Erro na requisicao:\n %s", bodyResponse)
        os.Exit(1)
    }
    defer response.Body.Close()

    var responseJson UserPlaylistsResponse
    if jsonErr := json.Unmarshal(bodyResponse, &responseJson); jsonErr != nil{
        log.Fatal(jsonErr)
    }

    playlistsUri = make(map[string]string)
    for _, userPlaylist := range responseJson.Items{
        playlistsUri[userPlaylist.Name] = userPlaylist.Tracks.Href
    }
    return
}

func getPlaylistTracks(accessToken string, playlistTracksURL string) (trackList []string){
    request, _ := http.NewRequest(http.MethodGet,
                                  playlistTracksURL,
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

    for _, item := range jsonResponse.Items{
        track := item.Track 
        trackFullName := fmt.Sprintf("%s;%s;%s\n", track.Name, track.Album.Artists[0].Name, track.Album.Name)
        trackList = append(trackList, trackFullName)
    }
    return
}

func SavePlaylistsTracks(accessToken string, references map[string]string){
    for name, ref := range references{
        trackList := getPlaylistTracks(accessToken, ref)
        fileName := fmt.Sprintf("playlists/%s.csv", name)
        trackFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
        if err != nil{
            log.Fatal(err)
        }
        for _, trackName := range trackList{
            trackFile.WriteString(trackName)
        }
    }
}
