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
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const fileName = "spotify/access_token.txt"

func generateBearerToken(spotifyClientId string, spotifyClientToken string) string{
    firstStep := fmt.Sprintf("%s:%s", spotifyClientId, spotifyClientToken)
    firstStepToBytes := []byte(firstStep)
    return base64.StdEncoding.EncodeToString(firstStepToBytes)
}


func GenerateNewAccessToken() string {
    enviroment := Env()

    formData := url.Values{}
    formData.Set("response_type", "code")
    formData.Set("scope", "user-library-read")
    formData.Set("client_id", enviroment.SpotifyClientId)
    formData.Set("redirect_uri", "http://localhost:8080/callback")
    formData.Set("state", enviroment.Secret)

    exec.Command("xdg-open", enviroment.SpotifyAccountsUrl+"/authorize?" + formData.Encode()).Start()

    chanToken := make(chan string)
    go func(){
        gin.DefaultWriter = io.Discard
        app := gin.New()
        app.GET("/callback", func(c *gin.Context){
            if enviroment.Secret != c.Query("state"){
                log.Fatal("Response from spotify malformed")
            }
            authToken := c.Query("code")
            chanToken <- authToken
            c.Abort()
        })
        app.Run(":8080")
    }()

    authorizationToken := <-chanToken
    formData = url.Values{
        "code": {authorizationToken},
        "redirect_uri": {"http://localhost:8080/callback"},
        "grant_type": {"authorization_code"},
    }
    authReq, err := http.NewRequest(
        http.MethodPost, 
        enviroment.SpotifyAccountsUrl+"/api/token", 
        strings.NewReader(formData.Encode()))
    if err != nil{
        log.Fatal(err)
    }
    authReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    stringToEncode := fmt.Sprintf("%s:%s", enviroment.SpotifyClientId, enviroment.SpotifyClientToken)
    authReq.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(stringToEncode)))
    authReq.Header.Set("Accept", "application/json")
    
    res, err := http.DefaultClient.Do(authReq)
    if err != nil || res.StatusCode != 200{
        log.Fatal(err)
    }
    defer res.Body.Close()

    bodyResponse, _ := io.ReadAll(res.Body)
    var responseJson ApiTokenResponse
    if err = json.Unmarshal(bodyResponse, &responseJson); err != nil{
        log.Fatal(err)
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
