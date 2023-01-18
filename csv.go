package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"spotify-crawler/google"

	"github.com/gocarina/gocsv"
)

func GetTracksFromCsv(playlistFile string) []google.Track{
    tracks := []google.Track{}
    trackFile, _ := os.OpenFile(playlistFile, os.O_RDONLY, os.ModePerm)

    gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
            r := csv.NewReader(in)
            r.Comma = ';'
            return r // Allows use pipe as delimiter
        })
    if err := gocsv.UnmarshalFile(trackFile, &tracks); err != nil{
        log.Fatalf("Erro no unmarshal: %s", err)
    }
    return tracks
}
