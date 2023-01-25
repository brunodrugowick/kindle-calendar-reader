package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/brunodrugowick/go-http-server-things/pkg/server"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
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

type displayEvent struct {
	Day         string
	TimeSlot    string
	Description string
}

func main() {
	web := server.
		NewDefaultServerBuilder().
		WithHandlerFunc("/", handlerRoot).
		Build()
	log.Fatal(web.ListenAndServe())
}

func getEvents() *calendar.Events {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := Bod(time.Now()).Format(time.RFC3339)
	log.Printf("Getting events starting at %v", t)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(100).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next hundred of the user's events: %v", err)
	}
	//fmt.Println("Upcoming events:")
	//if len(events.Items) == 0 {
	//	fmt.Println("No upcoming events found.")
	//} else {
	//	for _, item := range events.Items {
	//		date := item.Start.DateTime
	//		if date == "" {
	//			date = item.Start.Date
	//		}
	//		fmt.Printf("%v (%v)\n", item.Summary, date)
	//	}
	//}

	return events
}

// From https://stackoverflow.com/questions/25254443/return-local-beginning-of-day-time-object
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// From https://stackoverflow.com/questions/25254443/return-local-beginning-of-day-time-object
func Truncate(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	var displayEvents []displayEvent
	events := getEvents()
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
		displayEvents = append(displayEvents, displayEvent{
			Day:         day,
			TimeSlot:    timeSlot,
			Description: event.Summary,
		})
	}
	template, _ := template.ParseFiles("./templates/events-page.html")
	err := template.Execute(w, displayEvents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
