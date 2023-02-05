package json

import (
	"context"
	"encoding/json"
	"kindle-calendar-reader/pkg/api"
	"kindle-calendar-reader/pkg/service/events"
	"net/http"
)

type jsonApi struct {
	events events.Events
	path   string
}

func NewJsonApi(events events.Events, path string) api.Api {
	return &jsonApi{
		events: events,
		path:   path,
	}
}

func (api *jsonApi) HandleRequests(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")

	displayEvents, err := api.events.GetEvents(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(displayEvents)
	return
}

func (api *jsonApi) GetPath() string {
	return api.path
}
