package types

type GroupedEvents map[string][]DisplayEvent

func GroupEventsByDay(events []DisplayEvent) GroupedEvents {
	grouped := make(GroupedEvents)
	for _, event := range events {
		dateString := event.StartTimestamp.Format("2006-01-02")
		grouped[dateString] = append(grouped[dateString], event)
	}
	return grouped
}
