package bonusly

import (
	"fmt"
	"io"
	"net/http"
)

// New creates a new Client that can be used to interact with the Bonus.ly REST API.
func New(cfg Configuration, options ...ClientOption) *Client {
	c := &Client{
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

// ClientOption is a functional option to allow overwriting certain configuration options of the Client.
type ClientOption func(c *Client)

// WithHttpClient sets a new http.Client to be used by the bonusly.Client.
func WithHttpClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithEndpoint sets a new endpoint to be used by the bonusly.Client.
func WithEndpoint(endpoint Endpoint) ClientOption {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

// WithApplicationName sets a new application name to be used by the bonusly.Client.
//
// Setting the application name will set the via property of every bonus created by that application to the
// application name.
//
// See: https://bonusly.docs.apiary.io/#introduction/basics/your-application-name
func WithApplicationName(name string) ClientOption {
	return func(c *Client) {
		c.applicationName = name
	}
}

// Client is the main struct that implements all the methods for the different Bonus.ly REST API endpoints.
//
// The Client can be configured either through bonusly.ClientOption when using the bonusly.New() function or
// through the bonusly.Configuration.
//
// To use the client successfully, you need to have a Bonus.ly account and an API token.
//
// See: https://bonusly.docs.apiary.io/#
type Client struct {
	// httpClient is the client used for requests to the Bonus.ly REST API.
	//
	// To use your own client use bonusly.WithHttpClient option when creating a new bonusly.Client.
	httpClient *http.Client

	// token is the Bonus.ly authentication token that is used for requests to the Bonus.ly REST API.
	//
	// There are several kinds of tokens. Read-only tokens, Read and write tokens and admin tokens. All functions
	// require at least a read-only token, but some also require a write or admin token.
	//
	// The token can be set through the bonusly.Configuration.
	token string

	// endpoint is the HTTP endpoint the bonusly.Client sends requests to.
	//
	// By default, this is the production endpoint. The endpoint can be changed using the bonusly.WithEndpoint option
	// when creating a new client with bonusly.New().
	//
	// Default: EndpointProduction
	endpoint Endpoint

	// applicationName is the name of the application sending requests to the Bonus.ly REST API.
	//
	// The application name is stored in the 'via' property of certain resources to indicate which application created
	// or modified those resources. So choosing a meaningful name is important. For example: "HR Bot" or
	// "Account Report Application".
	//
	// The applicationName can be set using the bonusly.WithApplicationName option when creating a new bonusly.Client
	// using the bonusly.New() function.
	//
	// Default: DefaultApplicationName
	applicationName string
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("HTTP_APPLICATION_NAME", c.applicationName)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

// Endpoint is the Bonus.ly REST API endpoint to which requests are sent to.
type Endpoint string

const (
	// EndpointProduction is the fully qualified domain name and path of the Bonus.ly REST API production endpoint.
	//
	// To change the endpoint use the WithEndpoint option when creating a new Client using the New() function.
	EndpointProduction Endpoint = "https://bonus.ly/api/v1"
)

// DefaultApplicationName is the default application name used to send requests to the Bonus.ly REST API.
//
// To change the application name use the WithApplicationName option when creating a new Client using the
// New() function.
const DefaultApplicationName string = "Bonus.ly Go SDK"

// Configuration represents the Client configuration for the Bonus.ly API.
type Configuration struct {
	// Token is the access token used to authenticate requests with the Bonus.ly REST API.
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
