package auth

import (
	"context"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type Auth interface {
	GetConfiguredHttpClient(ctx context.Context) (*http.Client, error)
	GetRedirectUrl(host string) string
	GetTokenFromCode(ctx context.Context, authCode string)
}

type auth struct {
	GoogleAppConfig   *oauth2.Config
	autoRefreshClient *http.Client
}

func NewAuthService(oauthConfig *oauth2.Config) Auth {
	return &auth{
		GoogleAppConfig: oauthConfig,
	}
}

func (as *auth) GetConfiguredHttpClient(ctx context.Context) (*http.Client, error) {
	return as.autoRefreshClient, nil
}

func (as *auth) GetRedirectUrl(host string) string {
	authURL := as.GoogleAppConfig.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline)
	log.Printf("Redirect URL: %v", authURL)

	return authURL
}

func (as *auth) GetTokenFromCode(ctx context.Context, authCode string) {
	tok, err := as.GoogleAppConfig.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Printf("Unable to retrieve token from web: %v", err)
		return
	}
	as.autoRefreshClient = as.GoogleAppConfig.Client(ctx, tok)
}
