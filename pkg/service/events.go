package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"kindle-calendar-reader/pkg/types"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Service interface {
	GetEvents(ctx context.Context, token string) *calendar.Events
}

type V1 struct {
	GoogleAppConfig *oauth2.Config
	googleClient    *http.Client
}

const (
	defaultMaxEvents    int64  = 20
	defaultCalendarName string = "primary"
	defaultOrderBy      string = "startTime"
)

func NewEventsService(oauthConfig *oauth2.Config) V1 {
	return V1{
		GoogleAppConfig: oauthConfig,
		googleClient:    getGoogleConfiguredHTTPClient(oauthConfig),
	}
}

func (service *V1) GetEvents(ctx context.Context) []types.DisplayEvent {
	maxEvents := defaultMaxEvents
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(service.googleClient))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := truncateToStartOfDay(time.Now()).Format(time.RFC3339)
	log.Printf("Getting events starting at %v", t)
	events, err := srv.Events.List(defaultCalendarName).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(maxEvents).OrderBy(defaultOrderBy).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next %d of the user's events: %v", maxEvents, err)
	}
	var displayEvents []types.DisplayEvent
	for _, event := range events.Items {
		var day, timeSlot string
		if event.Start.DateTime != "" {
			day = strings.Split(event.Start.DateTime, "T")[0]
			timeSlot = strings.Split(strings.SplitAfter(event.Start.DateTime, "T")[1], "-")[0] +
				" - " +
				strings.Split(strings.SplitAfterN(event.End.DateTime, "T", 2)[1], "-")[0]

		} else {
			timeSlot = event.Start.Date
		}
		displayEvents = append(displayEvents, types.DisplayEvent{
			Day:         day,
			TimeSlot:    timeSlot,
			Description: event.Summary,
		})
	}

	return displayEvents
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

// Retrieve a token, saves the token, then returns the generated client.
func getGoogleConfiguredHTTPClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
