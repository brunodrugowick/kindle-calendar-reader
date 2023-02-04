package events

import (
	"context"
	"kindle-calendar-reader/pkg/api/types"
	"log"
)

type Events interface {
	GetEvents(ctx context.Context) ([]types.DisplayEvent, error)
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

func (delegator *eventsDelegator) GetEvents(ctx context.Context) (allEvents []types.DisplayEvent, err error) {
	for _, delegate := range delegator.delegates {
		events, err := delegate.GetEvents(ctx)
		if err != nil {
			log.Printf("Error getting events from delegator %s: %v", delegate, err)
			continue
		}
		allEvents = append(allEvents, events...)
	}
	return
}
