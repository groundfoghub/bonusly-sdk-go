package bonusly

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BaseUser struct {
	Id   string `json:"id"`
	Path string `json:"path"`

	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FullName     string `json:"full_name"`
	ShortName    string `json:"short_name"`
	DisplayName  string `json:"display_name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	ManagerEmail string `json:"manager_email"`

	FullPictureURL    *url.URL `json:"full_pic_url"`
	ProfilePictureURL *url.URL `json:"profile_pic_url"`

	Status       string    `json:"status"`
	IsAdmin      bool      `json:"admin"`
	LastActiveAt time.Time `json:"last_active_at"`
	CreatedAt    time.Time `json:"created_at"`
	HiredOn      time.Time `json:"hired_on,omitempty"`

	ExternalUniqueId string   `json:"external_unique_id"`
	BudgetBoost      int      `json:"budget_boost"`
	UserMode         UserMode `json:"user_mode"`
	Country          string   `json:"country"`
	TimeZone         string   `json:"time_zone"`

	CanReceive bool `json:"can_receive"`
	CanGive    bool `json:"can_give"`

	GiveAmounts          []int                 `json:"give_amounts"`
	SuggestedGiveAmounts []SuggestedGiveAmount `json:"suggested_give_amounts"`
	CustomProperties     CustomProperties      `json:"custom_properties"`

	ClientIds  ClientIDs `json:"client_ids"`
	IntercomId string    `json:"intercom_id"`
}

type User struct {
	BaseUser
}

type SuggestedGiveAmount struct {
	Amount int    `json:"dataValue"`
	Name   string `json:"name"`
}

type ClientIDs struct {
	Slack              string `json:"slack"`
	SlackId            string `json:"slack_id"`
	SlackDisplayName   string `json:"slack_display_name"`
	SlackHomeChannelID string `json:"slack_home_channel_id"`
}

type CustomProperties struct {
	Department string `json:"department"`
	Location   string `json:"location"`
	Role       string `json:"role"`
}

type UserMode string

const (
	UserModeNormal     UserMode = "normal"
	UserModeObserver   UserMode = "observer"
	UserModeReceiver   UserMode = "receiver"
	UserModeBenefactor UserMode = "benefactor"
	UserModeBot        UserMode = "bot"
)

type YYYYMMDD time.Time

func (u *BaseUser) UnmarshalJSON(data []byte) error {
	type Alias BaseUser

	user := struct {
		*Alias
		HiredOn           YYYYMMDD `json:"hired_on"`
		FullPictureURL    string   `json:"full_pic_url"`
		ProfilePictureURL string   `json:"profile_pic_url"`
	}{
		Alias: (*Alias)(u),
	}

	err := json.Unmarshal(data, &user)
	if err != nil {
		return err
	}

	user.Alias.HiredOn = time.Time(user.HiredOn)

	fullPicture, err := url.Parse(user.FullPictureURL)
	if err != nil {
		return err
	}
	user.Alias.FullPictureURL = fullPicture

	profilePicture, err := url.Parse(user.ProfilePictureURL)
	if err != nil {
		return err
	}
	user.Alias.ProfilePictureURL = profilePicture

	*u = BaseUser(*user.Alias)
	return nil
}

func (d *YYYYMMDD) UnmarshalJSON(b []byte) error {
	var hd string
	err := json.Unmarshal(b, &hd)
	if err != nil {
		return err
	}

	if hd == "" {
		return nil
	}

	parsed, err := time.Parse("2006-01-02", hd)
	if err != nil {
		return err
	}

	*d = YYYYMMDD(parsed)
	return nil
}

type SortProperty string

const (
	SortPropertyCreatedAt    SortProperty = "created_at"
	SortPropertyLastActiveAt SortProperty = "last_active_at"
	SortPropertyDisplayName  SortProperty = "display_name"
	SortPropertyFirstName    SortProperty = "first_name"
	SortPropertyLastName     SortProperty = "last_name"
	SortPropertyEmail        SortProperty = "email"
	SortPropertyCountry      SortProperty = "country"
	SortPropertyTimeZone     SortProperty = "time_zone"
)

type SortOrder int

const (
	SortOrderAscending SortOrder = iota
	SortOrderDescending
)

type ListUsersInput struct {
	Limit              int
	Skip               int
	Email              string
	CustomPropertyName string // TODO: Could be custom type that is easier to use properly
	SortBy             SortProperty
	SortOrder          SortOrder
	IncludeArchived    bool
	ShowFinancialData  bool
	UserMode           UserMode
}

type ListUsersOutput struct {
	Users []User
}

type ListUsersPaginatorClient interface {
	ListUsers(context.Context, *ListUsersInput) (*ListUsersOutput, error)
}

type ListUsersPaginator struct {
	client        ListUsersPaginatorClient
	params        *ListUsersInput
	firstPage     bool
	offset        int
	lastUserCount int
}

func NewListUsersPaginator(client ListUsersPaginatorClient, params *ListUsersInput) *ListUsersPaginator {
	if params == nil {
		params = &ListUsersInput{Limit: 20}
	}

	return &ListUsersPaginator{
		client:    client,
		params:    params,
		firstPage: true,
	}
}

func (p *ListUsersPaginator) HasMorePages() bool {
	return p.firstPage || p.lastUserCount >= p.params.Limit
}

func (p *ListUsersPaginator) NextPage(ctx context.Context) (*ListUsersOutput, error) {
	p.firstPage = false
	p.params.Skip = p.offset

	output, err := p.client.ListUsers(ctx, p.params)
	if err != nil {
		return nil, err
	}

	p.lastUserCount = len(output.Users)
	p.offset += p.lastUserCount
	return output, nil
}

func (c *Client) ListUsers(ctx context.Context, params *ListUsersInput) (*ListUsersOutput, error) {
	if params == nil {
		params = &ListUsersInput{}
	}

	u, err := newListUsersURL(c.endpoint, params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

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

		Result []User `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("list users: %v", r.Message)
	}

	return &ListUsersOutput{Users: r.Result}, nil
}

// newListUsersURL returns the URL to get a list of users (ListUsers) based on the provided endpoint and params. If the
// URL can not be created a non-nil error is returned and the URL is nil.
//nolint: cyclop
func newListUsersURL(endpoint Endpoint, params *ListUsersInput) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("%s/users", endpoint))
	if err != nil {
		return nil, err
	}

	q := u.Query()

	if params.Limit > 0 {
		q.Add("limit", strconv.Itoa(params.Limit))
	}

	if params.Skip > 0 {
		q.Add("skip", strconv.Itoa(params.Skip))
	}

	if params.Email != "" {
		q.Add("email", params.Email)
	}

	if params.CustomPropertyName != "" {
		q.Add("custom_property_name", params.CustomPropertyName)
	}

	if params.SortBy != "" {
		sortBy := ""
		if params.SortOrder == SortOrderDescending {
			sortBy += "-"
		}
		sortBy += string(params.SortBy)
		q.Add("sort", sortBy)
	}

	if params.IncludeArchived {
		q.Add("include_archived", "true")
	}

	if params.ShowFinancialData {
		q.Add("show_financial_data", "true")
	}

	if params.UserMode != "" {
		q.Add("user_mode", string(params.UserMode))
	}

	u.RawQuery = q.Encode()

	return u, nil
}

type GetUserInput struct {
	Id string
}

type GetUserOutput struct {
	User ExtendedUser
}

var (
	ErrMissingUserId = errors.New("missing user id")
)

type ExtendedUser struct {
	BaseUser

	EarningBalance               int    `json:"earning_balance"`
	EarningBalanceWithCurrency   string `json:"earning_balance_with_currency"`
	GiveBalance                  int    `json:"giving_balance"`
	GiveBalanceWithCurrency      string `json:"giving_balance_with_currency"`
	LifeTimeEarnings             int    `json:"lifetime_earnings"`
	LifeTimeEarningsWithCurrency string `json:"lifetime_earnings_with_currency"`
}

func (c *Client) GetUser(ctx context.Context, params *GetUserInput) (*GetUserOutput, error) {
	if params == nil {
		return nil, ErrMissingUserId
	}

	u := fmt.Sprintf("%s/users/%s", c.endpoint, params.Id)
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

		User ExtendedUser `json:"result"`
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return nil, fmt.Errorf("get user: %v", r.Message)
	}

	return &GetUserOutput{User: r.User}, nil
}
