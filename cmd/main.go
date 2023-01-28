package main

import (
	"github.com/brunodrugowick/go-http-server-things/pkg/server"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"kindle-calendar-reader/pkg/api"
	"kindle-calendar-reader/pkg/service"
	"log"
	"os"
	"strconv"
)

const defaultServerPort = 8080

func main() {

	googleAppConfig := setupGoogleAppClient()
	myService := service.NewEventsService(googleAppConfig)

	apiV1 := api.V1{Service: myService}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Printf("Cannot read server port from environment, setting default value: %d", defaultServerPort)
		serverPort = defaultServerPort
	}

	web := server.
		NewDefaultServerBuilder().
		SetPort(serverPort).
		WithHandlerFunc("/", apiV1.GetEvents).
		Build()

	log.Fatal(web.ListenAndServe())
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
