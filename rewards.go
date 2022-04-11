package bonusly

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ListRewardsInput struct {
	CatalogCountry string `json:"catalog_country"` // TODO: Reverse engineer and figure out if a "enum" can be used
	RequestCountry string `json:"request_country"` // TODO: Reverse engineer and figure out if a "enum" can be used
	PersonalizeFor string `json:"personalize_for"`
}

type RewardType string

const (
	RewardTypeUnknown   RewardType = "unknown"
	RewardTypeGiftCards RewardType = "gift_cards"
	RewardTypeDonations RewardType = "donations"
	RewardTypeCashOuts  RewardType = "cash_outs"
)

var rewardTypes map[string]RewardType = map[string]RewardType{
	"gift_cards": RewardTypeGiftCards,
	"donations":  RewardTypeDonations,
	"cash_outs":  RewardTypeCashOuts,
}

type RewardDenomination struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	DisplayPrice string `json:"display_price"`
}

type ListRewardsReward struct {
	Type                RewardType
	Name                string
	ImageUrl            *url.URL
	MinimumDisplayPrice string
	DescriptionText     string
	DescriptionHTML     string
	DisclaimerHTML      string
	Warning             string
	Categories          []string
	Denominations       []RewardDenomination
}

type ListRewardsOutput struct {
	Rewards []ListRewardsReward
}

type listRewardResponseResult struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Rewards []struct {
		Name                string `json:"name"`
		ImageUrl            string `json:"image_url"`
		MinimumDisplayPrice string `json:"minimum_display_price"`
		Description         struct {
			Text string `json:"text"`
			Html string `json:"html"`
		} `json:"description"`
		DisclaimerHtml string                                 `json:"disclaimer_html"`
		Warning        string                                 `json:"warning"`
		Categories     []string                               `json:"categories"`
		Denominations  []listRewardResponseResultDenomination `json:"denominations"`
	} `json:"rewards"`
}

type listRewardResponseResultDenomination struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	DisplayPrice string `json:"display_price"`
}

type listRewardsResponse struct {
	baseAPIResponse

	Result []listRewardResponseResult `json:"result"`
}

func (c *Client) ListRewards(ctx context.Context, params *ListRewardsInput) (*ListRewardsOutput, error) {
	if params == nil {
		params = &ListRewardsInput{}
	}

	u, err := url.Parse(fmt.Sprintf("%s/rewards", c.endpoint))
	if err != nil {
		return nil, err
	}

	q := u.Query()

	if params.CatalogCountry != "" {
		q.Add("catalog_country", params.CatalogCountry)
	}

	if params.RequestCountry != "" {
		q.Add("request_country", params.RequestCountry)
	}

	if params.PersonalizeFor != "" {
		q.Add("personalize_for", params.PersonalizeFor)
	}

	u.RawQuery = q.Encode()

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

	var r listRewardsResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("list rewards: %v", r.Message)
	}

	rewards := newRewards(r)

	return &ListRewardsOutput{Rewards: rewards}, nil
}

func newRewards(resp listRewardsResponse) []ListRewardsReward {
	rewards := make([]ListRewardsReward, 0)

	for i := range resp.Result {
		for _, r := range resp.Result[i].Rewards {
			iu, err := url.Parse(r.ImageUrl)
			if err != nil {
				// TODO: Do not just drop the error.
				continue
			}

			rewards = append(rewards, ListRewardsReward{
				Type:                newRewardType(resp.Result[i].Type),
				Name:                r.Name,
				ImageUrl:            iu,
				MinimumDisplayPrice: r.MinimumDisplayPrice,
				DescriptionText:     r.Description.Text,
				DescriptionHTML:     r.Description.Html,
				DisclaimerHTML:      r.DisclaimerHtml,
				Warning:             r.Warning,
				Categories:          r.Categories,
				Denominations:       newDenominations(r.Denominations),
			})
		}
	}

	return rewards
}

func newDenominations(denominations []listRewardResponseResultDenomination) []RewardDenomination {
	d := make([]RewardDenomination, len(denominations))

	for i := range denominations {
		d[i] = RewardDenomination{
			Id:           denominations[i].Id,
			Name:         denominations[i].Name,
			Price:        denominations[i].Price,
			DisplayPrice: denominations[i].DisplayPrice,
		}
	}

	return d
}

func newRewardType(t string) RewardType {
	rt, exists := rewardTypes[t]
	if !exists {
		return RewardTypeUnknown
	}

	return rt
}

type GetRewardInput struct {
	Id             string
	RequestCountry string
}

type GetRewardOutput struct {
	Reward GetRewardReward
}

type GetRewardReward struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Price        int      `json:"price"`
	DisplayPrice string   `json:"display_price"`
	Categories   []string `json:"categories"`
	Description  struct {
		Text string `json:"text"`
		Html string `json:"html"`
	} `json:"description"`
	ImageUrl       string `json:"image_url"`
	DisclaimerHtml string `json:"disclaimer_html"`
	Warning        string `json:"warning"`
	Quantity       int    `json:"quantity"`
	Type           string `json:"type"`
}

func (c *Client) GetReward(ctx context.Context, params *GetRewardInput) (*GetRewardOutput, error) {
	if params == nil {
		return nil, fmt.Errorf("user id missing")
	}

	u, err := url.Parse(fmt.Sprintf("%s/rewards/%s", c.endpoint, params.Id))
	if err != nil {
		return nil, err
	}

	q := u.Query()

	if params.RequestCountry != "" {
		q.Add("request_country", params.RequestCountry)
	}

	u.RawQuery = q.Encode()

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

		Result GetRewardReward `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("list users: %v", r.Message)
	}

	return &GetRewardOutput{Reward: r.Result}, nil
}
