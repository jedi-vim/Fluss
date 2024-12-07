package spotify

import (
	"encoding/json"
	"fluss/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)


func GetLikedArtists(finishedChan chan utils.ProgressStatus){
    env := Env()
    accessToken := GetAccessToken("user-follow-read")
    savedArtists := [][]string{
        {"name", "genres"},
    }

    continueFecthing := true
    after := ""
    for continueFecthing {
        queryParams := url.Values{}
        if after != ""{
            queryParams.Set("after", after)
        }
        queryParams.Set("limit", strconv.Itoa(50))
        queryParams.Set("type", "artist")
        request, _ := http.NewRequest(
            http.MethodGet,
            env.SpotifyApiUrl+"/v1/me/following?" + queryParams.Encode(),
            nil)
        request.Header.Add("Authorization", "Bearer " + accessToken)
        request.Header.Add("Content-Type", "application/json")
        response, err := http.DefaultClient.Do(request)
        if err != nil {
            log.Fatal(err)
        }

        bodyResponse, _ := io.ReadAll(response.Body)
        if response.StatusCode != 200{
            log.Fatalf("Erro durante a requisicao:\n %d %s", response.StatusCode, bodyResponse)
        }
        defer response.Body.Close()

        var jsonResponse ArtistsResponse
        if err := json.Unmarshal(bodyResponse, &jsonResponse); err != nil{
            log.Fatal(err)
        }

        if len(savedArtists) >= jsonResponse.Artists.Total{
            continueFecthing = false
            continue
        }

        for _, artist := range jsonResponse.Artists.Items {
            artistMap := []string{
                artist.Name,
                strings.Join(artist.Genres, "/"),
            }
            savedArtists = append(savedArtists, artistMap)
        }
        after = jsonResponse.Artists.Cursors.After
        progressStatus := utils.ProgressStatus{
            Finished: false,
            Message: fmt.Sprintf("Artistas recuperados %d de %d", len(savedArtists), jsonResponse.Artists.Total),
        }
        finishedChan <- progressStatus
    }
    csvProgressStatus := utils.ProgressStatus{
        Finished: false, 
        Message: fmt.Sprintf("Salvando Artistas em arquivo .csv"),
    }
    finishedChan <- csvProgressStatus
    utils.PersistCsv("spotify.liked-artists.csv", savedArtists)

    finalProgressStatus := utils.ProgressStatus{
        Finished: true,
        Message: fmt.Sprintf("Finalizado com total de Artistas recuperados %d\n", len(savedArtists)),
    }
    finishedChan <- finalProgressStatus
    return
}
