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

type events struct {
	autoRefreshClient *http.Client
	oauthConfig       *oauth2.Config
}

const (
	defaultMaxEvents           int64  = 20
	defaultCalendarName        string = "primary"
	defaultOrderBy             string = "startTime"
	timePortionOfRFC3339Format string = "T00:00:00+00:00"
)

func NewGoogleEventsService(auth *oauth2.Config) Events {
	return &events{
		oauthConfig: auth,
	}
}

func (service *events) GetRedirectUrl(host string) string {
	authURL := service.oauthConfig.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline)
	log.Printf("Redirect URL: %v", authURL)

	return authURL
}

func (service *events) GetTokenFromCode(ctx context.Context, authCode string) bool {
	tok, err := service.oauthConfig.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Printf("Unable to retrieve token from web: %v", err)
		return false
	}
	service.autoRefreshClient = service.oauthConfig.Client(ctx, tok)
	return true
}

func (service *events) GetProvider() string {
	return "Google"
}

func (service *events) Name() string {
	return "Google Service"
}

func (service *events) GetEventsStartingAt(ctx context.Context, start time.Time, limit int64) ([]types.DisplayEvent, error) {
	displayEvents, err := service.getEvents(ctx, start, limit)
	if err != nil {
		return []types.DisplayEvent{}, err
	}

	return displayEvents, nil
}

func (service *events) getEvents(ctx context.Context, startDate time.Time, limit int64) ([]types.DisplayEvent, error) {
	if limit < 1 {
		limit = defaultMaxEvents
	}
	var displayEvents []types.DisplayEvent
	client := service.autoRefreshClient

	// TODO start service when new'in this up?
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Calendar client: %v", err)
		return displayEvents, errors.New("unable to retrieve Calendar client")
	}

	log.Printf("Getting events starting at %v", startDate)
	googleEvents, err := srv.Events.
		List(defaultCalendarName).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(startDate.Format(time.RFC3339)).
		MaxResults(limit).
		OrderBy(defaultOrderBy).
		Do()
	if err != nil {
		log.Printf("Unable to retrieve next %d of the user's events: %v", limit, err)
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
