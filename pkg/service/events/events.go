package events

import (
	"context"
	"errors"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"kindle-calendar-reader/pkg/api/types"
	"kindle-calendar-reader/pkg/service/auth"
	"log"
	"strings"
	"time"
)

type Events interface {
	GetEvents(ctx context.Context) ([]types.DisplayEvent, error)
}

type events struct {
	authService auth.Service
}

const (
	defaultMaxEvents    int64  = 20
	defaultCalendarName string = "primary"
	defaultOrderBy      string = "startTime"
)

func NewEventsService(authService auth.Service) Events {
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

	t := truncateToStartOfDay(time.Now()).Format(time.RFC3339)
	log.Printf("Getting events starting at %v", t)
	maxEvents := defaultMaxEvents
	events, err := srv.Events.
		List(defaultCalendarName).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(maxEvents).
		OrderBy(defaultOrderBy).
		Do()
	if err != nil {
		log.Printf("Unable to retrieve next %d of the user's events: %v", maxEvents, err)
		return displayEvents, errors.New("error retrieving events")
	}

	for _, event := range events.Items {
		var day, timeSlot string
		if event.Start.DateTime != "" {
			day = strings.Split(event.Start.DateTime, "T")[0]
			timeSlot = shamelessCuttingThingsAway(event)
		} else {
			timeSlot = event.Start.Date
		}

		displayEvents = append(displayEvents, types.DisplayEvent{
			Day:         day,
			TimeSlot:    timeSlot,
			Description: event.Summary,
		})
	}

	return displayEvents, nil
}

func shamelessCuttingThingsAway(event *calendar.Event) string {
	startString := strings.Split(strings.SplitAfter(event.Start.DateTime, "T")[1], "-")[0]
	separator := " - "
	endString := strings.Split(strings.SplitAfterN(event.End.DateTime, "T", 2)[1], "-")[0]

	return startString + separator + endString
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
