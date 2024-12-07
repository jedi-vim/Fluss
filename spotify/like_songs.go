package spotify

import (
	"encoding/json"
	"fluss/utils"
	"io"
        "fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)


func GetLikedSongs(finishedChan chan utils.ProgressStatus){
    env := Env()
    accessToken := GetAccessToken("user-library-read")
    savedSongs := [][]string{
        {"artist", "album", "released_at", "song"},
    }

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
            trackMap := []string{
                track.Track.Artists[0].Name,
                track.Track.Album.Name,
                track.Track.Album.ReleaseDate,
                track.Track.Name,
            }
            savedSongs = append(savedSongs, trackMap)
        }
        offSet = offSet + 50
        progressStatus := utils.ProgressStatus{
            Finished: false,
            Message: fmt.Sprintf("Musicas recuperadas %d", len(savedSongs)),
        }
        finishedChan <- progressStatus
    }
    csvProgressStatus := utils.ProgressStatus{
        Finished: false, 
        Message: fmt.Sprintf("Salvando musicas em arquivo .csv"),
    }
    finishedChan <- csvProgressStatus
    utils.PersistCsv("spotify.liked-songs.csv", savedSongs)

    finalProgressStatus := utils.ProgressStatus{
        Finished: true, 
        Message: fmt.Sprintf("Finalizado com total de musicas recuperadas %d\n", len(savedSongs)),
    }
    finishedChan <- finalProgressStatus
    return
}
