package main

import (
    "os"
    "log"
    "time"
)

const fileName = "access_token.txt"

func getAccessTokenFromFile() string{
    log.Println("Acessando informacao do arquivo de token")
    fileInfo, err := os.Stat(fileName)
    if os.IsNotExist(err){
        log.Println("Arquivo de token inexistente.")
        return ""
    }
    if err != nil{
        log.Println("Um erro aconteceu e deu merda")
        log.Fatal(err)
    }
    lastMod := fileInfo.ModTime()
    tokenExpired := time.Since(lastMod).Hours() > float64(1)
    if tokenExpired{
        log.Println("Arquivo contem token ja expirado.")
        return ""
    }
    tokenBytes, err := os.ReadFile(fileName)
    if err != nil{
        log.Fatal("Erro ao ler o token")
    }
    return string(tokenBytes)
}

func saveToken(newAccessToken string){
    tokenFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
    if err != nil{
        log.Fatal(err)
    }
    tokenFile.WriteString(newAccessToken)
}
