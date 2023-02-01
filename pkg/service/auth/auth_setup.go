package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

type AuthSetup interface {
	GetConfiguredHttpClient(ctx context.Context) (*http.Client, error)
	GetRedirectUrl(host string) string
	GetTokenFromCode(ctx context.Context, authCode string)
}

type authSetup struct {
	GoogleAppConfig *oauth2.Config
}

const (
	tokFile string = "token.json"
)

func NewAuthSetupService(oauthConfig *oauth2.Config) AuthSetup {
	return &authSetup{
		GoogleAppConfig: oauthConfig,
	}
}

func (as *authSetup) GetConfiguredHttpClient(ctx context.Context) (*http.Client, error) {
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Printf("Could not read token file from %s", tokFile)
		return nil, errors.New("can't find user token")
	}
	return as.GoogleAppConfig.Client(ctx, tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (as *authSetup) GetRedirectUrl(host string) string {
	authURL := as.GoogleAppConfig.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("redirect_uri", fmt.Sprintf("http://%s/", host)))
	log.Printf("Redirect URL: %v", authURL)

	return authURL
}

func (as *authSetup) GetTokenFromCode(ctx context.Context, authCode string) {
	tok, err := as.GoogleAppConfig.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Printf("Unable to retrieve token from web: %v", err)
		return
	}
	saveToken(tokFile, tok)
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
