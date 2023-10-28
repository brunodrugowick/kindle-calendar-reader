package events

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"kindle-calendar-reader/pkg/api/types"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type googleEvents struct {
	abstractService
	client *http.Client
}

const (
	defaultMaxEvents           int64  = 20
	defaultCalendarName        string = "primary"
	defaultOrderBy             string = "startTime"
	timePortionOfRFC3339Format string = "T00:00:00+00:00"
)

func NewGoogleEventsService(auth *oauth2.Config) Events {
	return &googleEvents{
		abstractService: abstractService{
			oauthConfig:  auth,
			logger:       log.New(log.Writer(), "Google Service ", 3),
			providerName: "Google",
		},
	}
}

func (service *googleEvents) GetTokenFromCode(ctx context.Context, authCode string) bool {
	tok, err := service.oauthConfig.Exchange(context.TODO(), authCode)
	if err != nil {
		service.logger.Printf("Unable to retrieve token from web: %v", err)
		return false
	}
	service.client = service.oauthConfig.Client(ctx, tok)
	return true
}

func (service *googleEvents) GetEventsStartingAt(ctx context.Context, start time.Time, limit int64) ([]types.DisplayEvent, error) {
	displayEvents, err := service.getEvents(ctx, start, limit)
	if err != nil {
		return []types.DisplayEvent{}, err
	}

	return displayEvents, nil
}

func (service *googleEvents) getEvents(ctx context.Context, startDate time.Time, limit int64) ([]types.DisplayEvent, error) {
	if limit < 1 {
		limit = defaultMaxEvents
	}
	var displayEvents []types.DisplayEvent
	client := service.client

	// TODO start service when new'in this up?
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		service.logger.Printf("Unable to retrieve Calendar client: %v", err)
		return displayEvents, errors.New("unable to retrieve Calendar client")
	}

	service.logger.Printf("Getting events starting at %v", startDate)
	googleEvents, err := srv.Events.
		List(defaultCalendarName).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(startDate.Format(time.RFC3339)).
		MaxResults(limit).
		OrderBy(defaultOrderBy).
		Do()
	if err != nil {
		service.logger.Printf("Unable to retrieve next %d of the user's events: %v", limit, err)
		return displayEvents, errors.New("error retrieving events from Google")
	}

	for _, event := range googleEvents.Items {
		var start, end time.Time
		var allDay bool
		if event.Start.DateTime == "" { // All day events only have the .Start.Date value
			start, err = time.Parse(time.RFC3339, event.Start.Date+timePortionOfRFC3339Format)
			allDay = true
		} else {
			start, err = time.Parse(time.RFC3339, event.Start.DateTime)
			if err != nil {
				continue
			}

			end, err = time.Parse(time.RFC3339, event.End.DateTime)
			allDay = false
		}

		displayEvents = append(displayEvents, types.DisplayEvent{
			Day:            fmt.Sprintf("%s %02d", start.Month().String(), start.Day()),
			StartTime:      fmt.Sprintf("%02d:%02d", start.Hour(), start.Minute()),
			EndTime:        fmt.Sprintf("%02d:%02d", end.Hour(), end.Minute()),
			StartTimestamp: start,
			AllDay:         allDay,
			Description:    event.Summary,
		})
	}

	return displayEvents, nil
}
