package utils

import (
	"encoding/csv"
	"log"
	"os"
)

func PersistCsv(fileName string, sequence [][]string){
    fileCsv, err := os.Create(fileName)
    if err != nil{
        log.Fatal(err)
    }
    defer fileCsv.Close()

    writer := csv.NewWriter(fileCsv)
    defer writer.Flush()

    if err = writer.WriteAll(sequence);err != nil{
        log.Fatal("Erro ao salvar conteudo em .csv")
    }
}
