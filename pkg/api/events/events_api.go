package events

import (
	"context"
	"html/template"
	"kindle-calendar-reader/pkg/api"
	"kindle-calendar-reader/pkg/api/types"
	"kindle-calendar-reader/pkg/service/events"
	"log"
	"net/http"
)

const todayPageTemplate = `<html><title>Today</title><body>
{{range $date, $events := .}}
	<h2>{{$date}}</h2>
	<ul>
		{{range $events}}
			<u>{{if .AllDay}}{{"All Day"}}{{else}}{{.StartTime}} - {{.EndTime}}{{end}}</u> {{.Description}}<br>
		{{end}}
	</ul>
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
	displayEvents, err := api.events.GetEventsStartingToday(ctx, 0)
	if err != nil {
		http.Redirect(w, r, "/setup", http.StatusFound)
	}

	groupEventsByDay := types.GroupEventsByDay(displayEvents)

	tmpl, err := template.New("Today").Parse(todayPageTemplate)
	if err != nil {
		log.Printf("Error creating template: %v", err)
	}
	err = tmpl.Execute(w, groupEventsByDay)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
