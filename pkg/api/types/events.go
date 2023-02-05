package types

import "time"

type DisplayEvent struct {
	Day            string    `json:"day"`
	StartTime      string    `json:"startTime"`
	StartTimestamp time.Time `json:"startTimestamp"`
	EndTime        string    `json:"endTime"`
	AllDay         bool      `json:"allDay"`
	Description    string    `json:"description"`
}
