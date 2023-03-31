package events

import (
	"context"
	"errors"
	"fmt"
	"kindle-calendar-reader/pkg/api/types"
	"kindle-calendar-reader/pkg/service/auth"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type events struct {
	authService auth.Auth
}

const (
	defaultMaxEvents           int64  = 20
	defaultCalendarName        string = "primary"
	defaultOrderBy             string = "startTime"
	timePortionOfRFC3339Format string = "T15:04:05Z07:00"
)

func NewGoogleEventsService(authService auth.Auth) Events {
	return &events{
		authService: authService,
	}
}

func (service *events) Name() string {
	return "Google Service"
}

func (service *events) GetEventsStartingToday(ctx context.Context, limit int64) ([]types.DisplayEvent, error) {
	timeMin := startOfDay(time.Now())
	displayEvents, err := service.getEvents(ctx, timeMin, limit)
	if err != nil {
		return []types.DisplayEvent{}, err
	}

	return displayEvents, nil
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
	client, err := service.authService.GetConfiguredHttpClient(ctx)
	if err != nil {
		log.Printf("Could not get a configured HTTP client due to err: %v", err)
		return displayEvents, fmt.Errorf("could not get events: %w", err)
	}
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

// endOfDay is from https://stackoverflow.com/questions/25254443/return-local-beginning-of-day-time-object
func endOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// startOfDay is from https://stackoverflow.com/questions/25254443/return-local-beginning-of-day-time-object
func startOfDay(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}
