package api

import (
	"context"
	"html/template"
	"kindle-calendar-reader/pkg/api/types"
	"kindle-calendar-reader/pkg/service/auth"
	"kindle-calendar-reader/pkg/service/events"
	"log"
	"net/http"
)

type Api interface {
	DispatchRootRequests(w http.ResponseWriter, r *http.Request)
	DispatchSetupRequests(w http.ResponseWriter, r *http.Request)
}

type api struct {
	Events events.Events
	Auth   auth.AuthSetup
}

const todayPageTemplate = `<html><body>
<h1>Today</h1>
{{range .}}
	<div class="container">
		<p><small><b>{{.Day}}</b> {{.TimeSlot}}</small><br>{{.Description}}
	</div>
{{end}}
</body></html>`

const setupApiTokenTemplate = `<html><body>
<div class="container">
	<p><a href={{.GoogleUrl}}>Setup Google</a>
</div>
</body></html>`

func NewEventsApi(events events.Events, auth auth.AuthSetup) Api {
	return &api{
		Events: events,
		Auth:   auth,
	}
}

func (api *api) DispatchRootRequests(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	switch len(code) {
	case 0:
		getEvents(ctx, w, r, api)
	default:
		api.Auth.GetTokenFromCode(ctx, code)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

}

func (api *api) DispatchSetupRequests(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodGet:
		setupTokenGetRequest(w, r, api)
	case http.MethodPost:
		setupTokenPostRequest(ctx, w, r, api)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getEvents(ctx context.Context, w http.ResponseWriter, r *http.Request, api *api) {
	client, err := api.Auth.GetConfiguredHttpClient(ctx)
	if err != nil {
		redirectURL := "/setup"
		log.Printf("Could not get a configured HTTP client. Redirecting user to %s", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}

	displayEvents, err := api.Events.GetEvents(ctx, client)
	if err != nil {
		http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
	}

	tmpl, err := template.New("Today").Parse(todayPageTemplate)
	if err != nil {
		log.Printf("Error creating template: %v", err)
	}
	err = tmpl.Execute(w, displayEvents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func setupTokenGetRequest(w http.ResponseWriter, r *http.Request, api *api) {
	displayToken := types.DisplaySetupInfo{
		GoogleUrl: api.Auth.GetRedirectUrl(r.Host),
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

func setupTokenPostRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, api *api) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.Form.Get("Token")
	api.Auth.GetTokenFromCode(ctx, code)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
