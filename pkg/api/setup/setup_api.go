package setup

import (
	"context"
	"html/template"
	"kindle-calendar-reader/pkg/api"
	"kindle-calendar-reader/pkg/service/events"
	"log"
	"net/http"
)

const setupApiTokenTemplate = `<html>
<head><style>
.dark-mode {
  background-color: #121212; /* Dark background color */
  color: #aaaaaa; /* Light text color */
}
</style></head>
<body class="dark-mode">
<div class="container">
	<h1> Available providers </h1>
	{{ range $key, $value := . }}
		{{if eq $key "Outlook"}}
			<p><a href="{{$value}}">Setup {{$key}} (you need access to the app logs in its current state to properly sign-in to Outlook)</a>
		{{else}}
			<p><a href="{{$value}}">Setup {{$key}}</a>
		{{end}}
	{{end}}
</div>

<h2>Pages</h2>

- <a href="/">Home</a>

</body></html>`

type setupApi struct {
	eventsServices []events.Events
	path           string
}

func NewSetupApi(path string, eventsServices ...events.Events) api.Api {
	return &setupApi{
		eventsServices: eventsServices,
		path:           path,
	}
}

const (
	codeQueryParam     = "code"
	providerQueryParam = "state"
)

func (a *setupApi) HandleRequests(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	queryParams, err := api.ParseFormAndGetFromRequest(r, codeQueryParam, providerQueryParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch len(queryParams[codeQueryParam]) {
	case 0:
		setupRouteGetRequest(w, r, a)
	default:
		for _, provider := range a.eventsServices {
			if provider.GetProviderName() == queryParams[providerQueryParam] &&
				provider.GetTokenFromCode(ctx, queryParams[codeQueryParam]) {
				http.Redirect(w, r, "/", http.StatusFound)
			}
		}
	}
}

func (a *setupApi) GetPath() string {
	return a.path
}

func setupRouteGetRequest(w http.ResponseWriter, r *http.Request, api *setupApi) {
	providerTokens := make(map[string]string)
	for _, provider := range api.eventsServices {
		providerTokens[provider.GetProviderName()] = provider.GetRedirectUrl()
	}

	tmpl, err := template.New("Setup").Parse(setupApiTokenTemplate)
	if err != nil {
		log.Printf("Error creating template: %v", err)
		return
	}

	err = tmpl.Execute(w, providerTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
