package json

import (
	"context"
	"encoding/json"
	"kindle-calendar-reader/pkg/api"
	"kindle-calendar-reader/pkg/api/types"
	"kindle-calendar-reader/pkg/service/events"
	"net/http"
	"strconv"
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

const startDateQueryParam = "startDate"
const limitQueryParam = "limit"

func (a *jsonApi) HandleRequests(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	queryParams, queryParamsErr := api.ParseFormAndGetFromRequest(r, startDateQueryParam, limitQueryParam)

	limit, queryParamsErr := strconv.ParseInt(queryParams[limitQueryParam], 10, 64)
	if queryParamsErr != nil {
		limit = 20
	}

	var displayEvents []types.DisplayEvent
	var serviceErr error
	switch len(queryParams[startDateQueryParam]) {
	case 0:
		displayEvents, serviceErr = a.events.GetEventsStartingToday(ctx, limit)
	default:
		startTime, conversionErr := time.Parse(time.RFC3339, queryParams[startDateQueryParam])
		if conversionErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		displayEvents, serviceErr = a.events.GetEventsStartingAt(ctx, startTime, limit)
	}
	if serviceErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(displayEvents)
	return
}

func (a *jsonApi) GetPath() string {
	return a.path
}
