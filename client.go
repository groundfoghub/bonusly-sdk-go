package bonusly

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// New creates a new bonusly client that can be used to interact with the bonusly API.
func New(cfg Configuration, options ...ClientOption) Client {
	c := &apiClient{
		httpClient:      http.DefaultClient,
		token:           cfg.Token,
		endpoint:        EndpointProduction,
		applicationName: DefaultApplicationName,
	}

	for _, fn := range options {
		fn(c)
	}

	return c
}

type ClientOption func(c *apiClient)

func WithHttpClient(client *http.Client) ClientOption {
	return func(c *apiClient) {
		c.httpClient = client
	}
}

func WithEndpoint(endpoint Endpoint) ClientOption {
	return func(c *apiClient) {
		c.endpoint = endpoint
	}
}

func WithApplicationName(name string) ClientOption {
	return func(c *apiClient) {
		c.applicationName = name
	}
}

type Client interface {
	GetUser(context.Context, *GetUserInput) (*GetUserOutput, error)
	ListUsers(context.Context, *ListUsersInput) (*ListUsersOutput, error)
	CreateBonus(context.Context, *CreateBonusInput) (*CreateBonusOutput, error)
	GetRedemption(context.Context, *GetRedemptionInput) (*GetRedemptionOutput, error)
	ListRedemptions(context.Context, *ListRedemptionsInput) (*ListRedemptionsOutput, error)
	GetReward(context.Context, *GetRewardInput) (*GetRewardOutput, error)
	ListRewards(context.Context, *ListRewardsInput) (*ListRewardsOutput, error)
	ListWebhooks(context.Context) (*ListWebhooksOutput, error)
	CreateWebhook(context.Context, *CreateWebhookInput) (*CreateWebhookOutput, error)
	UpdateWebhook(context.Context, *UpdateWebhookInput) (*UpdateWebhookOutput, error)
	DeleteWebhook(context.Context, *DeleteWebhookInput) (*DeleteWebhookOutput, error)
}

type apiClient struct {
	httpClient      *http.Client
	token           string
	endpoint        Endpoint
	applicationName string
}

func (c *apiClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("HTTP_APPLICATION_NAME", c.applicationName)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

type Endpoint string

const (
	EndpointProduction Endpoint = "https://bonus.ly/api/v1"
)

const DefaultApplicationName string = "Bonus.ly Go SDK"

type Configuration struct {
	Token string
}

func closeCloser(c io.Closer) {
	err := c.Close()
	if err != nil {
		panic(err)
	}
}

type baseAPIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
