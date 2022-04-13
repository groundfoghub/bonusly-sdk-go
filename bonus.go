package bonusly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CreateBonusInput struct {
	GiverEmail    string
	Receivers     []string
	Reason        string
	Amount        uint
	ParentBonusID string
}

type createBonusBody struct {
	GiverEmail    string `json:"giver_email"`
	Reason        string `json:"reason"`
	ParentBonusID string `json:"parent_bonus_id,omitempty"`
}

type CreateBonusOutput struct {
}

type createBonusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func (c *Client) CreateBonus(ctx context.Context, params *CreateBonusInput) (*CreateBonusOutput, error) {
	// TODO: Add option that allows overwriting the formatting of the reason string.
	b := createBonusBody{
		GiverEmail:    params.GiverEmail,
		Reason:        newReason(params),
		ParentBonusID: params.ParentBonusID,
	}

	body, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/bonuses", c.endpoint), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := readAndCloseBody(resp)
	if err != nil {
		return nil, err
	}

	var r createBonusResponse
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("create bonus: %s", r.Message)
	}

	return &CreateBonusOutput{}, nil
}

func newReason(params *CreateBonusInput) string {
	// string builder is used because the receivers list could potentially contain >50 users.
	receivers := strings.Builder{}
	for i := range params.Receivers {
		receivers.WriteString("@")
		receivers.WriteString(strings.TrimSpace(params.Receivers[i]))
		receivers.WriteString(" ")
	}

	return fmt.Sprintf("+%d %s %s", params.Amount, strings.TrimSpace(receivers.String()), params.Reason)
}
