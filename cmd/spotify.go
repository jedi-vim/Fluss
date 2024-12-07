package cmd

import (
	"encoding/csv"
	"fluss/spotify"
	"log"
        "os"
	"github.com/spf13/cobra"
)

func PersistCsv(mapSequence []map[string]string){
    fileCsv, err := os.Create("spotify.liked-songs.csv")
    if err == nil{
        log.Fatal(err)
    }
    defer fileCsv.Close()

    writer := csv.NewWriter(fileCsv)
    defer writer.Flush()

    for idx, elem := range mapSequence{
        if idx == 0 {
            var keys []string = []string{}
            for k, _ :=  range elem{
                keys = append(keys, k)
            }
            writer.Write(keys)
        }
        var elemValues []string = []string{}
        for _, v := range elem{
            elemValues = append(elemValues, v)
        }
        writer.Write(elemValues)
    }
}

var spotifyCmd = &cobra.Command{
        Use: "spotify <content-to-migrate>",
        Short: "Export spotify content [liked-music, artists]",
        Args: cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string){
            accessToken := spotify.GetAccessToken()
            switch args[0]{
                case "liked-music":
                    log.Printf("Solicitacao: exportar: [%s], origem: [Spotify]\n", args[0])
                    likedSongs := spotify.GetLikedSongs(accessToken)
                    PersistCsv(likedSongs)
                    log.Printf("Total de musicas recuperadas %d", len(likedSongs))
                default:
                    log.Printf("Opcao nao existe: %s", args[0])
            }
        },
    }

func init(){
    RootCmd.AddCommand(spotifyCmd)
}
