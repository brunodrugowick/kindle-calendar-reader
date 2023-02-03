package events

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"kindle-calendar-reader/pkg/api/types"
	"kindle-calendar-reader/pkg/service/auth"
	"log"
	"time"
)

type Events interface {
	GetEvents(ctx context.Context) ([]types.DisplayEvent, error)
}

type events struct {
	authService auth.Auth
}

const (
	defaultMaxEvents           int64  = 20
	defaultCalendarName        string = "primary"
	defaultOrderBy             string = "startTime"
	timePortionOfRFC3339Format string = "T15:04:05Z07:00"
)

func NewEventsService(authService auth.Auth) Events {
	return &events{
		authService: authService,
	}
}

func (service *events) GetEvents(ctx context.Context) ([]types.DisplayEvent, error) {
	var displayEvents []types.DisplayEvent

	client, err := service.authService.GetConfiguredHttpClient(ctx)
	if err != nil {
		log.Printf("Could not get a configured HTTP client due to err: %v", err)
		return displayEvents, errors.New("could not get events")
	}
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Calendar client: %v", err)
		return displayEvents, errors.New("unable to retrieve Calendar client")
	}

	now := time.Now()
	startOfTheDay := truncateToStartOfDay(now).Format(time.RFC3339)
	endOfTheDay := truncateToEndOfDay(now).Format(time.RFC3339)

	log.Printf("Getting events starting at %v", startOfTheDay)
	maxEvents := defaultMaxEvents
	events, err := srv.Events.
		List(defaultCalendarName).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(startOfTheDay).
		TimeMax(endOfTheDay).
		MaxResults(maxEvents).
		OrderBy(defaultOrderBy).
		Do()
	if err != nil {
		log.Printf("Unable to retrieve next %d of the user's events: %v", maxEvents, err)
		return displayEvents, errors.New("error retrieving events")
	}

	for _, event := range events.Items {
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
			Day:         fmt.Sprintf("%s %02d", start.Month().String(), start.Day()),
			StartTime:   fmt.Sprintf("%02d:%02d", start.Hour(), start.Minute()),
			EndTime:     fmt.Sprintf("%02d:%02d", end.Hour(), end.Minute()),
			AllDay:      allDay,
			Description: event.Summary,
		})
	}

	return displayEvents, nil
}

// truncateToStartOfDay is from https://stackoverflow.com/questions/25254443/return-local-beginning-of-day-time-object
func truncateToStartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// truncateToEndOfDay is from https://stackoverflow.com/questions/25254443/return-local-beginning-of-day-time-object
func truncateToEndOfDay(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}
