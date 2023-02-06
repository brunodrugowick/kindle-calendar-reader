package json

import (
	"context"
	"encoding/json"
	"kindle-calendar-reader/pkg/api"
	"kindle-calendar-reader/pkg/api/types"
	"kindle-calendar-reader/pkg/service/events"
	"net/http"
	"time"
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

func (a *jsonApi) HandleRequests(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")

	startString, err := api.ParseFormAndGetFromRequest(r, "startDate")
	var displayEvents []types.DisplayEvent
	switch len(startString) {
	case 0:
		displayEvents, err = a.events.GetEventsStartingToday(ctx)
	default:
		startTime, err := time.Parse(time.RFC3339, startString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		displayEvents, err = a.events.GetEventsStartingAt(ctx, startTime)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(displayEvents)
	return
}

func (a *jsonApi) GetPath() string {
	return a.path
}
