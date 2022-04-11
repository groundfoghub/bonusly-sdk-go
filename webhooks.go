package bonusly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// WebhookEventType represents the different event types a webhook can be subscribed to.
type WebhookEventType string

const (
	WebhookEventTypeBonusCreated            WebhookEventType = "bonus.created"
	WebhookEventTypeAchievementEventCreated WebhookEventType = "achievement_event.created"
)

// Webhook represents a Bonus.ly webhook. A webhook has a unique identifier and a URL that is called if one of the
// subscribed event types is triggered.
type Webhook struct {
	// ID of the webhook.
	ID string `json:"id"`
	// URL of the webhook where messages by Bonus.ly are POST'ed to.
	URL *url.URL `json:"url"`
	// EventTypes represents the list of events to be notified of.
	EventTypes []WebhookEventType `json:"event_types"`
}

// UnmarshalJSON is a custom json.Unmarshaler for the Webhook type to properly deserialize the Webhook.URL.
func (w *Webhook) UnmarshalJSON(data []byte) error {
	type Alias Webhook

	webhook := struct {
		*Alias
		URL string `json:"url"`
	}{
		Alias: (*Alias)(w),
	}

	err := json.Unmarshal(data, &webhook)
	if err != nil {
		return err
	}

	u, err := url.Parse(webhook.URL)
	if err != nil {
		return err
	}

	webhook.Alias.URL = u

	*w = Webhook(*webhook.Alias)
	return nil
}

// ListWebhooksOutput represents the output of the "List Webhooks" operation.
type ListWebhooksOutput struct {
	// Webhooks is a slice of all found webhooks. If no webhooks are found the slice will be empty.
	Webhooks []Webhook
}

// ListWebhooks returns all webhooks.
//
// Note: The Bonus.ly API does not support pagination for this API. Therefore, no paginator exists.
func (c *Client) ListWebhooks(ctx context.Context) (*ListWebhooksOutput, error) {
	u := fmt.Sprintf("%s/webhooks", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeCloser(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type response struct {
		baseAPIResponse

		Result []Webhook `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("list webhooks: %v", r.Message)
	}

	return &ListWebhooksOutput{Webhooks: r.Result}, nil
}

// CreateWebhookInput represents the input of the "Create Webhook" operation.
type CreateWebhookInput struct {
	// URL of the webhook where messages by Bonus.ly are POST'ed to.
	URL *url.URL `json:"url"`
	// EventTypes represents the list of events to be notified of.
	EventTypes []WebhookEventType `json:"event_types"`
}

// CreateWebhookOutput represents the output of the "Create Webhook" operation.
type CreateWebhookOutput struct {
	// ID of the newly created webhook.
	ID string `json:"id"`
}

// CreateWebhook creates a new webhook.
func (c *Client) CreateWebhook(ctx context.Context, params *CreateWebhookInput) (*CreateWebhookOutput, error) {
	if params == nil {
		return nil, fmt.Errorf("params missing")
	}

	reqBody := struct {
		URL        string             `json:"url"`
		EventTypes []WebhookEventType `json:"event_types"`
	}{
		URL:        params.URL.String(),
		EventTypes: params.EventTypes,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/webhooks", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeCloser(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type response struct {
		baseAPIResponse

		Result struct {
			Id string `json:"id"`
		} `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("create webhook: %v", r.Message)
	}

	return &CreateWebhookOutput{ID: r.Result.Id}, nil
}

// UpdateWebhookInput represents the input of the "Update Webhook" operation.
type UpdateWebhookInput struct {
	// ID of the webhook to update.
	ID string `json:"id"`
	// URL of the webhook where messages by Bonus.ly are POST'ed to.
	URL *url.URL `json:"url"`
	// EventTypes represents the list of events to be notified of (optional).
	EventTypes []WebhookEventType `json:"event_types"`
}

// MarshalJSON is a custom json.Marshaler for the UpdateWebhookInput type to properly serialize the UpdateWebhookInput.URL.
func (w *UpdateWebhookInput) MarshalJSON() ([]byte, error) {
	var u *string
	if w.URL != nil {
		t := w.URL.String()
		u = &t
	}

	return json.Marshal(&struct {
		ID         string             `json:"id"`
		URL        *string            `json:"url,omitempty"`
		EventTypes []WebhookEventType `json:"event_types"`
	}{
		ID:         w.ID,
		URL:        u,
		EventTypes: w.EventTypes,
	})
}

// UpdateWebhookOutput represents the output of the "Update Webhook" operation.
type UpdateWebhookOutput struct {
	// ID of the updated webhook.
	ID string `json:"id"`
}

// UpdateWebhook updates a single webhook.
func (c *Client) UpdateWebhook(ctx context.Context, params *UpdateWebhookInput) (*UpdateWebhookOutput, error) {
	if params == nil {
		return nil, fmt.Errorf("params missing")
	}

	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/webhooks/%s", c.endpoint, params.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeCloser(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type response struct {
		baseAPIResponse

		Result struct {
			Id string `json:"id"`
		} `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("create webhook: %v", r.Message)
	}

	return &UpdateWebhookOutput{ID: r.Result.Id}, nil
}

// DeleteWebhookInput represents the input of the "Delete Webhook" operation.
type DeleteWebhookInput struct {
	// ID of the webhook to delete.
	ID string `json:"id"`
}

// DeleteWebhookOutput represents the output of the "Delete Webhook" operation.
type DeleteWebhookOutput struct {
	// ID is the ID of the deleted webhook.
	ID string `json:"id"`
}

// DeleteWebhook deletes a webhook with the provided id.
func (c *Client) DeleteWebhook(ctx context.Context, params *DeleteWebhookInput) (*DeleteWebhookOutput, error) {
	if params == nil {
		return nil, fmt.Errorf("params missing")
	}

	u := fmt.Sprintf("%s/webhooks/%s", c.endpoint, params.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeCloser(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type response struct {
		baseAPIResponse

		Result struct {
			Id string `json:"id"`
		} `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("create webhook: %v", r.Message)
	}

	return &DeleteWebhookOutput{ID: r.Result.Id}, nil
}
