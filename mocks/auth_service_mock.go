package mocks

import (
	"context"
	"net/http"
)

type AuthService struct {
	Mocker
}

func (a AuthService) GetConfiguredHttpClient(ctx context.Context) (*http.Client, error) {
	returns := a.called("GetConfiguredHttpClient", ctx)
	return returns.Get(0).(*http.Client), returns.Error(1)
}

func (a AuthService) GetRedirectUrl(host string) string {
	returns := a.called("GetRedirectUrl", host)
	return returns.Get(0).(string)
}

func (a AuthService) GetTokenFromCode(ctx context.Context, authCode string) {
	a.called("GetRedirectUrl", ctx, authCode)
}
