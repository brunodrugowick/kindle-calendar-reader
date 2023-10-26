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
	autoRefreshClient *msgraphsdk.GraphServiceClient
	oauthConfig       *oauth2.Config
}

func NewOutlookEventsService(auth *oauth2.Config) Events {
	return &outlookEvents{
		oauthConfig: auth,
	}
}

func (service *outlookEvents) GetRedirectUrl() string {
	authURL := service.oauthConfig.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline)
	log.Printf("Redirect URL: %v", authURL)

	return authURL
}

func (service *outlookEvents) GetTokenFromCode(ctx context.Context, authCode string) bool {

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

	result, err := client.Me().Get(ctx, nil)
	if err != nil {
		fmt.Printf("Error doing stuff: %v", err)
	}
	fmt.Printf("%v", *result.GetId())

	service.autoRefreshClient = client
	return true
}

func (service *outlookEvents) GetProviderName() string {
	return "Outlook"
}

func (service *outlookEvents) GetEventsStartingAt(ctx context.Context, start time.Time, limit int64) ([]types.DisplayEvent, error) {
	displayEvents, err := service.getEvents(ctx, start, limit)
	if err != nil {
		return []types.DisplayEvent{}, err
	}

	return displayEvents, nil
}

func (service *outlookEvents) getEvents(ctx context.Context, startDate time.Time, limit int64) ([]types.DisplayEvent, error) {
	if limit < 1 {
		limit = defaultMaxEvents
	}
	var displayEvents []types.DisplayEvent
	_ = service.autoRefreshClient

	//res, err := client.Get("https://graph.microsoft.com/v1.0/users/me/calendar/events")
	//if err != nil {
	//	log.Printf("Unable to retrieve events: %v", err)
	//	return displayEvents, errors.New("unable to retrieve Calendar events")
	//}

	//log.Printf(res.Status)
	return displayEvents, nil
}
