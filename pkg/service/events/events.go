package events

import (
	"context"
	"kindle-calendar-reader/pkg/api/types"
	"log"
	"time"
)

type Events interface {
	GetEventsStartingToday(ctx context.Context) ([]types.DisplayEvent, error)
	GetEventsStartingAt(ctx context.Context, time time.Time) ([]types.DisplayEvent, error)
}

type eventsDelegator struct {
	delegates []Events
}

func NewEventsDelegator(eventService ...Events) Events {
	var delegator eventsDelegator
	for i := 0; i < len(eventService); i++ {
		delegator.delegates = append(delegator.delegates, eventService[i])
	}
	return &delegator
}

func (delegator *eventsDelegator) GetEventsStartingToday(ctx context.Context) (allEvents []types.DisplayEvent, err error) {
	for _, delegate := range delegator.delegates {
		events, err := delegate.GetEventsStartingToday(ctx)
		if err != nil {
			log.Printf("Error getting events from delegator %s: %v", delegate, err)
			continue
		}
		allEvents = append(allEvents, events...)
	}
	return
}

func (delegator *eventsDelegator) GetEventsStartingAt(ctx context.Context, start time.Time) (allEvents []types.DisplayEvent, err error) {
	for _, delegate := range delegator.delegates {
		events, err := delegate.GetEventsStartingAt(ctx, start)
		if err != nil {
			log.Printf("Error getting events from delegator %s: %v", delegate, err)
			continue
		}
		allEvents = append(allEvents, events...)
	}
	return
}
