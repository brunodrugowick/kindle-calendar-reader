package main

import (
	"github.com/brunodrugowick/go-http-server-things/pkg/server"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"kindle-calendar-reader/pkg/api"
	eventsApi "kindle-calendar-reader/pkg/api/events"
	"kindle-calendar-reader/pkg/api/json"
	"kindle-calendar-reader/pkg/api/setup"
	"kindle-calendar-reader/pkg/service/auth"
	"kindle-calendar-reader/pkg/service/events"
	"log"
	"os"
	"strconv"
)

const defaultServerPort = 8080

func main() {

	// Services
	googleAppConfig := setupGoogleAppClient()
	authService := auth.NewAuthService(googleAppConfig)
	googleEventsService := events.NewGoogleEventsService(authService)
	eventsService := events.NewEventsDelegator(googleEventsService)

	// APIs
	var apis []api.Api
	apis = append(apis, eventsApi.NewEventsApi(eventsService, "/"))
	apis = append(apis, setup.NewSetupApi(authService, "/setup"))
	apis = append(apis, json.NewJsonApi(eventsService, "/json"))

	// Server
	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Printf("Cannot read server port from environment, setting default value: %d", defaultServerPort)
		serverPort = defaultServerPort
	}
	serverBuilder := server.
		NewDefaultServerBuilder().
		SetPort(serverPort)
	for _, a := range apis {
		serverBuilder.WithHandlerFunc(a.GetPath(), a.HandleRequests)
	}
	srv := serverBuilder.Build()
	log.Fatal(srv.ListenAndServe())
}

func setupGoogleAppClient() *oauth2.Config {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}
