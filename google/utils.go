package google

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/api/youtube/v3"
)

type Track struct{
    Song   string        `csv:"song"`
    Artist string        `csv:"artist"`
    Album  string        `csv:"album"`
}

func CreatePlaylist(name string, environ *Settings) *youtube.Playlist{
    youtubeService := GetYoutubeService(environ)
    insertCall := youtubeService.Playlists.Insert([]string{"snippet", "status"}, 
        &youtube.Playlist{
            Snippet: &youtube.PlaylistSnippet{
                Title: name,
            },
            Status: &youtube.PlaylistStatus{
                PrivacyStatus: "public",
            },
        },
    )
    playlistData, err := insertCall.Do()
    if err != nil{
        log.Fatal(err)
    }
    return playlistData 
} 

func PopulatePlaylist(playlist *youtube.Playlist, tracks []string, environ *Settings){
    fmt.Printf("Adicionando Musicas a playlist %s...\n", playlist.Snippet.Title)
    ytService := GetYoutubeService(environ)
    var partParam = []string{"snippet", "id"}
    for _, track := range tracks{
        searchResponse, err := ytService.Search.List(partParam).Q(track).MaxResults(3).Do()
        checkError(err)

        var mostRelevant *youtube.SearchResult = searchResponse.Items[0]
        time.Sleep(3 * time.Second)
        fmt.Printf("Adicionando: (%s) -> (%s)\t", track, mostRelevant.Snippet.Title)
        addedItem, err := ytService.PlaylistItems.Insert(partParam,
            &youtube.PlaylistItem{
                Snippet: &youtube.PlaylistItemSnippet{
                    PlaylistId: playlist.Id,
                    ResourceId: &youtube.ResourceId{
                        VideoId: mostRelevant.Id.VideoId,
                        Kind: "youtube#video",
                    },
                },
            }).Do()

        checkError(err)
        fmt.Printf("%s\n", addedItem.Id)
        time.Sleep(5 * time.Second)
    }
}
