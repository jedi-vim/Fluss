package spotify

import (
    "os"
    "io"
    "log"
    "time"
    "fmt"
    "strings"
    "net/http"
    "net/url"
    "encoding/base64"
    "encoding/json"
)

const fileName = "spotify/access_token.txt"

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

func getAccessTokenFromFile() string{
    log.Println("Acessando informacao do arquivo de token")
    fileInfo, err := os.Stat(fileName)
    if os.IsNotExist(err){
        log.Println("Arquivo de token inexistente.")
        return ""
    }
    if err != nil{
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
        log.Fatalf("Erro ao abrir arquivo de token: %s", err)
    }
    tokenFile.WriteString(newAccessToken)
}

func GetAccessToken() (accessToken string){
    accessToken = getAccessTokenFromFile() 
    if accessToken == ""{
        log.Println("No valid token exists, a new one will be generated")
        accessToken = GenerateNewAccessToken()
        saveToken(accessToken)
    }
    return
}
