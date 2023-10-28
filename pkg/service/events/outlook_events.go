package events

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"golang.org/x/oauth2"
	"kindle-calendar-reader/pkg/api/types"
	"log"
	"time"
)

type outlookEvents struct {
	abstractService
	client *msgraphsdk.GraphServiceClient
}

func NewOutlookEventsService(auth *oauth2.Config) Events {
	return &outlookEvents{
		abstractService: abstractService{
			oauthConfig:  auth,
			logger:       log.New(log.Writer(), "Outlook Service ", 3),
			providerName: "Outlook",
		},
	}
}

func (service *outlookEvents) GetTokenFromCode(ctx context.Context, authCode string) bool {
	// TODO Well, that's how it's going to be for now!
	cred, err := azidentity.NewDeviceCodeCredential(
		&azidentity.DeviceCodeCredentialOptions{
			AdditionallyAllowedTenants: []string{"9efaa0ad-666e-48e3-9405-f981ae695b78"},
			ClientID:                   service.oauthConfig.ClientID,
			DisableInstanceDiscovery:   false,
			TenantID:                   "9efaa0ad-666e-48e3-9405-f981ae695b78",
			UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
				fmt.Println(message.Message)
				return nil
			},
		},
	)
	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
	}

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, service.oauthConfig.Scopes)
	if err != nil {
		service.logger.Printf("Error doing stuff: %v", err)
		return false
	}

	service.client = client
	return true
}

func (service *outlookEvents) GetEventsStartingAt(ctx context.Context, start time.Time, limit int64) ([]types.DisplayEvent, error) {
	if service.client == nil {
		return []types.DisplayEvent{}, nil
	}

	displayEvents, err := service.getEvents(ctx, start, limit)
	if err != nil {
		return displayEvents, err
	}

	return displayEvents, nil
}

func (service *outlookEvents) getEvents(ctx context.Context, startDate time.Time, limit int64) ([]types.DisplayEvent, error) {
	var displayEvents []types.DisplayEvent
	client := service.client

	res, err := client.Me().Events().Get(ctx, nil)
	if err != nil {
		service.logger.Printf("Could not get events: {}", err)
		return displayEvents, err
	}

	// TODO Need to make it actually work and then...
	// TODO Extract events into displayEvents

	service.logger.Printf(fmt.Sprintf("%d", *res.GetOdataCount()))
	return displayEvents, err
}
