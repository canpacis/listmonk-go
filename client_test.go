package listmonkgo_test

import (
	"context"
	"os"
	"testing"

	listmonkgo "github.com/canpacis/listmonk-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	godotenv.Load()
	code := m.Run()
	os.Exit(code)
}

func createClient() *listmonkgo.Client {
	return listmonkgo.New(
		listmonkgo.WithBaseURL(os.Getenv("LISTMONK_URL")),
		listmonkgo.WithAPIUser(os.Getenv("API_USER")),
		listmonkgo.WithToken(os.Getenv("API_TOKEN")),
	)
}

func TestSubscriberEndpoints(t *testing.T) {
	assert := assert.New(t)
	client := createClient()

	var err error
	ctx := context.Background()

	createParams := &listmonkgo.CreateSubscriberParams{
		Email:                   "test@example.com",
		Name:                    "Test User",
		Status:                  listmonkgo.EnabledSubscriberStatus,
		Lists:                   []int{},
		Attributes:              map[string]any{},
		PreconfirmSubscriptions: true,
	}
	subscriber, err := client.CreateSubscriber(ctx, createParams)
	assert.NoError(err)
	assert.Equal(createParams.Email, subscriber.Email)
	assert.Equal(createParams.Name, subscriber.Name)
	assert.Equal(createParams.Status, subscriber.Status)
	assert.Equal(len(createParams.Lists), len(subscriber.Lists))

	list, err := client.CreateList(ctx, &listmonkgo.CreateListParams{
		Name:        "Test",
		Type:        listmonkgo.PublicTypeListEntry,
		Optin:       listmonkgo.SingleOptinListEntry,
		Tags:        []string{},
		Description: "Subscriber Endpoints Test",
	})
	assert.NoError(err)

	created, err := client.CreateSubscription(ctx, &listmonkgo.CreateSubscriptionParams{
		Email:     subscriber.Email,
		Name:      subscriber.Name,
		ListUUIDs: []uuid.UUID{list.UUID},
	})
	assert.NoError(err)
	assert.Equal(true, created)
	deleted, err := client.DeleteList(ctx, list.ID)
	assert.NoError(err)
	assert.Equal(true, deleted)

	deleted, err = client.DeleteSubscriber(ctx, subscriber.ID)
	assert.NoError(err)
	assert.Equal(true, deleted)
}

func TestListEndpoints(t *testing.T) {
	assert := assert.New(t)
	client := createClient()

	var err error
	ctx := context.Background()

	_, err = client.GetLists(ctx, nil)
	assert.NoError(err)

	_, err = client.GetPublicLists(ctx)
	assert.NoError(err)

	createParams := &listmonkgo.CreateListParams{
		Name:        "Test",
		Type:        listmonkgo.PrivateTypeListEntry,
		Optin:       listmonkgo.SingleOptinListEntry,
		Tags:        []string{"test tag"},
		Description: "Test description",
	}
	list, err := client.CreateList(ctx, createParams)
	assert.NoError(err)
	assert.Equal(createParams.Name, list.Name)
	assert.Equal(createParams.Description, list.Description)
	assert.Equal(createParams.Type, list.Type)
	assert.Equal(createParams.Optin, list.Optin)
	assert.Equal(len(createParams.Tags), len(list.Tags))

	updateParams := &listmonkgo.UpdateListParams{
		Name:        "New Name",
		Description: "New description",
	}
	newList, err := client.UpdateList(ctx, list.ID, updateParams)
	assert.NoError(err)
	assert.Equal(updateParams.Name, newList.Name)
	assert.Equal(updateParams.Description, newList.Description)

	deleted, err := client.DeleteList(ctx, list.ID)
	assert.NoError(err)
	assert.Equal(true, deleted)
}

func TestImportEndpoints(t *testing.T) {
	assert := assert.New(t)
	client := createClient()

	ctx := context.Background()
	logs, err := client.GetImportLogs(ctx)
	assert.NoError(err)
	assert.Equal("", logs)
}

func TestTransactionalMail(t *testing.T) {
	assert := assert.New(t)
	client := createClient()

	ctx := context.Background()
	sent, err := client.SendTemplate(ctx, &listmonkgo.SendTemplateParams{
		SubscriberEmail: "canpacis@gmail.com",
		TemplateID:      6,
		FromEmail:       "Onboarding <onboarding@canpacis.com>",
		Data: map[string]any{
			"CTAText": "Confirm Your Email",
			"CTAURL":  "https://www.canpacis.com",
		},
	})
	assert.NoError(err)
	assert.Equal(true, sent)
}
