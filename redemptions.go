package bonusly

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Redemption struct {
	Id             string   `json:"id"`
	UserId         string   `json:"user_id"`
	UserEmail      string   `json:"user_email"`
	GifteeEmail    string   `json:"giftee_email"`
	Title          string   `json:"title"`
	AmountInPoints int      `json:"amount_in_points"`
	AmountInUsd    string   `json:"amount_in_usd"`
	Categories     []string `json:"categories"`
}

type ListRedemptionsInput struct {
	Limit int
	Skip  int
}

type ListRedemptionsOutput struct {
	Redemptions []Redemption
}

type ListRedemptionsPaginatorClient interface {
	ListRedemptions(context.Context, *ListRedemptionsInput) (*ListRedemptionsOutput, error)
}

type ListRedemptionsPaginator struct {
	client          ListRedemptionsPaginatorClient
	params          *ListRedemptionsInput
	firstPage       bool
	offset          int
	lastResultCount int
}

func NewListRedemptionsPaginator(client ListRedemptionsPaginatorClient, params *ListRedemptionsInput) *ListRedemptionsPaginator {
	if params == nil {
		params = &ListRedemptionsInput{
			Limit: 100,
			Skip:  0,
		}
	}

	return &ListRedemptionsPaginator{
		client:    client,
		params:    params,
		firstPage: true,
	}
}

func (p *ListRedemptionsPaginator) HasMorePages() bool {
	return p.firstPage || p.lastResultCount >= p.params.Limit
}

func (p *ListRedemptionsPaginator) NextPage(ctx context.Context) (*ListRedemptionsOutput, error) {
	p.firstPage = false
	p.params.Skip = p.offset

	output, err := p.client.ListRedemptions(ctx, p.params)
	if err != nil {
		return nil, err
	}

	p.lastResultCount = len(output.Redemptions)
	p.offset += p.lastResultCount

	return output, nil
}

func (c *Client) ListRedemptions(ctx context.Context, params *ListRedemptionsInput) (*ListRedemptionsOutput, error) {
	if params == nil {
		params = &ListRedemptionsInput{
			Limit: 100,
			Skip:  0,
		}
	}

	u, err := url.Parse(fmt.Sprintf("%s/redemptions", c.endpoint))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("limit", strconv.Itoa(params.Limit))
	q.Add("skip", strconv.Itoa(params.Skip))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

		Result []Redemption `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("list redemptions: %v", r.Message)
	}

	return &ListRedemptionsOutput{Redemptions: r.Result}, nil
}

type GetRedemptionInput struct {
	Id string
}

type GetRedemptionOutput struct {
	Redemption GetRedemptionRedemption
}

type GetRedemptionRedemption struct {
	Id             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	State          string    `json:"state"`
	CertificateUrl string    `json:"certificate_url"`
	ClaimUrl       string    `json:"claim_url"`
	AutoApprovable bool      `json:"auto_approvable"`
	RewardDetails  struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Price        int    `json:"price"`
		DisplayPrice string `json:"display_price"`
		Type         string `json:"type"`
		ImageUrl     string `json:"image_url"`
	} `json:"reward_details"`
}

func (c *Client) GetRedemption(ctx context.Context, params *GetRedemptionInput) (*GetRedemptionOutput, error) {
	if params == nil {
		return nil, fmt.Errorf("id missing")
	}

	u, err := url.Parse(fmt.Sprintf("%s/redemptions/%s", c.endpoint, params.Id))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

		Result GetRedemptionRedemption `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("get redemption: %v", r.Message)
	}

	return &GetRedemptionOutput{Redemption: r.Result}, nil
}
