package events

import (
	"context"
	"kindle-calendar-reader/pkg/api/types"
	"log"
	"time"
)

type Events interface {
	GetRedirectUrl(host string) string
	GetTokenFromCode(ctx context.Context, authCode string) bool
	GetProvider() string
	GetEventsStartingToday(ctx context.Context, limit int64) ([]types.DisplayEvent, error)
	GetEventsStartingAt(ctx context.Context, time time.Time, limit int64) ([]types.DisplayEvent, error)
	Name() string
}

type Delegator interface {
	GetEventsStartingToday(ctx context.Context, limit int64) (allEvents []types.DisplayEvent, err error)
	GetEventsStartingAt(ctx context.Context, start time.Time, limit int64) (allEvents []types.DisplayEvent, err error)
}

type eventsDelegator struct {
	delegates []Events
}

func NewEventsDelegator(eventService ...Events) Delegator {
	var delegator eventsDelegator
	for i := 0; i < len(eventService); i++ {
		delegator.delegates = append(delegator.delegates, eventService[i])
	}
	return &delegator
}

func (delegator *eventsDelegator) GetEventsStartingToday(ctx context.Context, limit int64) (allEvents []types.DisplayEvent, err error) {
	for _, delegate := range delegator.delegates {
		events, err := delegate.GetEventsStartingToday(ctx, limit)
		if err != nil {
			log.Printf("Error getting events from delegator %s: %v", delegate.Name(), err)
			continue
		}
		allEvents = append(allEvents, events...)
	}
	return
}

func (delegator *eventsDelegator) GetEventsStartingAt(ctx context.Context, start time.Time, limit int64) (allEvents []types.DisplayEvent, err error) {
	for _, delegate := range delegator.delegates {
		events, err := delegate.GetEventsStartingAt(ctx, start, limit)
		if err != nil {
			log.Printf("Error getting events from delegator %s: %v", delegate, err)
			continue
		}
		allEvents = append(allEvents, events...)
	}
	return
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
