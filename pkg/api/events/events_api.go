package events

import (
	"context"
	"html/template"
	"kindle-calendar-reader/pkg/api"
	"kindle-calendar-reader/pkg/service/events"
	"log"
	"net/http"
)

const todayPageTemplate = `<html><title>Today</title><body>
{{range .}}
	<div class="container">
		<p><small>{{if .AllDay}} <b>All Day</b> {{else}}<b>{{.Day}}</b> {{.StartTime}} - {{.EndTime}} {{end}}</small><br>{{.Description}}
	</div>
{{end}}
<p>That's it for today! See you tomorrow.
</body></html>`

type eventsApi struct {
	events events.Events
	path   string
}

func NewEventsApi(events events.Events, path string) api.Api {
	return &eventsApi{
		events: events,
		path:   path,
	}
}

func (api *eventsApi) HandleRequests(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	getEvents(ctx, w, r, api)
}

func (api *eventsApi) GetPath() string {
	return api.path
}

func getEvents(ctx context.Context, w http.ResponseWriter, r *http.Request, api *eventsApi) {
	displayEvents, err := api.events.GetEvents(ctx)
	if err != nil {
		http.Redirect(w, r, "/setup", http.StatusFound)
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
