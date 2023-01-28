package api

import (
	"context"
	"google.golang.org/appengine/log"
	"html/template"
	"kindle-calendar-reader/pkg/service"
	"net/http"
)

type Api interface {
	getEvents()
	setupAccount()
}

type V1 struct {
	Service service.V1
}

const todayPageTemplate = `<html><body>
<h1>Today</h1>
{{range .}}<p><small><b>{{.Day}}</b> {{.TimeSlot}}</small><br>{{.Description}}{{end}}
</body></html>`

func (api *V1) GetEvents(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	displayEvents := api.Service.GetEvents(ctx)
	tmpl, err := template.New("Today").Parse(todayPageTemplate)
	if err != nil {
		log.Infof(ctx, "Error creating template: %v", err)
	}
	err = tmpl.Execute(w, displayEvents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
