package types

type DisplayEvent struct {
	Day         string `json:"day"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	AllDay      bool   `json:"allDay"`
	Description string `json:"description"`
}
