package cmd

import (
	"fluss/spotify"
	"fluss/utils"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var spotifyCmd = &cobra.Command{
        Use: "spotify <content-to-migrate>",
        Short: "Export spotify content [liked-music, artists]",
        Args: cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string){
            beginMessage := fmt.Sprintf("Solicitacao: exportar: [%s], origem: [Spotify]\n", args[0])
            switch args[0]{
                case "liked-music":
                    log.Printf(beginMessage)
                    utils.AnimationBar(spotify.GetLikedSongs)
                case "liked-artists":
                    log.Printf(beginMessage)
                    utils.AnimationBar(spotify.GetLikedArtists)
                default:
                    log.Printf("Opcao nao existe: %s", args[0])
            }
        },
    }

func init(){
    RootCmd.AddCommand(spotifyCmd)
}
