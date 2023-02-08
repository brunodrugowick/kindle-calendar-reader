package setup

import (
	"context"
	"html/template"
	"kindle-calendar-reader/pkg/api"
	"kindle-calendar-reader/pkg/api/types"
	"kindle-calendar-reader/pkg/service/auth"
	"log"
	"net/http"
)

const setupApiTokenTemplate = `<html><body>
<div class="container">
	<p><a href={{.GoogleUrl}}>Setup Google</a>
</div>
</body></html>`

type setupApi struct {
	auth auth.Auth
	path string
}

func NewSetupApi(authService auth.Auth, path string) api.Api {
	return &setupApi{
		auth: authService,
		path: path,
	}
}

const codeQueryParam = "code"

func (a *setupApi) HandleRequests(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	queryParams, err := api.ParseFormAndGetFromRequest(r, codeQueryParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch len(queryParams[codeQueryParam]) {
	case 0:
		setupRouteGetRequest(w, r, a)
	default:
		a.auth.GetTokenFromCode(ctx, queryParams[codeQueryParam])
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (a *setupApi) GetPath() string {
	return a.path
}

func setupRouteGetRequest(w http.ResponseWriter, r *http.Request, api *setupApi) {
	displayToken := types.DisplaySetupInfo{
		GoogleUrl: api.auth.GetRedirectUrl(r.Host),
	}

	tmpl, err := template.New("Setup").Parse(setupApiTokenTemplate)
	if err != nil {
		log.Printf("Error creating template: %v", err)
		return
	}

	err = tmpl.Execute(w, displayToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
