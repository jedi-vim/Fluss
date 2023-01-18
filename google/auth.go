package google

import (
    "os"
    "fmt"
    "context"
    "log"
    "net/http"

    "google.golang.org/api/youtube/v3"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/option"
)

func checkError(err error){
    if err != nil{
        log.Fatal(err)
    }
}

func getClient(config *oauth2.Config, tokFile string) *http.Client {
        // The file token.json stores the user's access and refresh tokens, and is
        // created automatically when the authorization flow completes for the first
        // time.
        tok, err := tokenFromFile(tokFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(tokFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web %v", err)
        }
        return tok
}

func GetYoutubeService(environ *Settings) (*youtube.Service){
    cxt := context.Background()
    
    credentialJson, err := os.ReadFile(environ.GoogleOAuthCredentials)
    checkError(err)

    config, err := google.ConfigFromJSON(
        credentialJson, 
        youtube.YoutubeScope,
    )
    checkError(err)
    client := getClient(config, environ.GoogleTokenFile)

    youtubeService, err := youtube.NewService(cxt, 
        option.WithHTTPClient(client))
    checkError(err)
    return youtubeService
}
