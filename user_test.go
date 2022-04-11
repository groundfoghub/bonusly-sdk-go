package bonusly

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestUser_UnmarshalHiredOn(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			"ok",
			args{data: []byte(`{"hired_on": "2022-03-01"}`)},
			time.Date(2022, 03, 01, 0, 0, 0, 0, time.UTC),
			false,
		},
		{
			"no-hire-date",
			args{data: []byte(`{}`)},
			time.Time{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got User
			err := json.Unmarshal(tt.args.data, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got.HiredOn != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, Want %v", got, tt.want)
			}
		})
	}
}

func TestUser_UnmarshalFullPictureURL(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			"ok",
			args{data: []byte(`{"full_pic_url": "https://bonusly.s3.amazonaws.com/uploads/user/default_avatar/123456789/full_avatar.png"}`)},
			mustURL(t, "https://bonusly.s3.amazonaws.com/uploads/user/default_avatar/123456789/full_avatar.png"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got User
			err := json.Unmarshal(tt.args.data, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if *got.FullPictureURL != *tt.want {
				t.Errorf("UnmarshalJSON() got = %v, Want %v", got.FullPictureURL, tt.want)
			}
		})
	}
}

func TestUser_UnmarshalProfilePictureURL(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			"ok",
			args{data: []byte(`{"profile_pic_url": "https://bonusly.s3.amazonaws.com/uploads/user/default_avatar/123456789/profile_avatar.png"}`)},
			mustURL(t, "https://bonusly.s3.amazonaws.com/uploads/user/default_avatar/123456789/profile_avatar.png"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got User
			err := json.Unmarshal(tt.args.data, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if *got.ProfilePictureURL != *tt.want {
				t.Errorf("UnmarshalJSON() got = %v, Want %v", got.FullPictureURL, tt.want)
			}
		})
	}
}

func TestUser_UnmarshalUserMode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    UserMode
		wantErr bool
	}{
		{
			"normal",
			args{data: []byte(`{"user_mode": "normal"}`)},
			UserModeNormal,
			false,
		},
		{
			"observer",
			args{data: []byte(`{"user_mode": "observer"}`)},
			UserModeObserver,
			false,
		},
		{
			"receiver",
			args{data: []byte(`{"user_mode": "receiver"}`)},
			UserModeReceiver,
			false,
		},
		{
			"benefactor",
			args{data: []byte(`{"user_mode": "benefactor"}`)},
			UserModeBenefactor,
			false,
		},
		{
			"bot",
			args{data: []byte(`{"user_mode": "bot"}`)},
			UserModeBot,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got User
			err := json.Unmarshal(tt.args.data, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got.UserMode != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, Want %v", got.FullPictureURL, tt.want)
			}
		})
	}
}

type mockClient struct {
	pages []*ListUsersOutput
	err   error
	c     int
}

func (m *mockClient) ListUsers(context.Context, *ListUsersInput) (*ListUsersOutput, error) {
	p := m.pages[m.c]
	m.c++

	return p, m.err
}

func TestListUsersPaginator(t *testing.T) {
	client := &mockClient{
		pages: []*ListUsersOutput{
			{Users: newUsers(t, 20)},
			{Users: newUsers(t, 20)},
			{Users: newUsers(t, 5)},
		},
	}

	var users []User

	paginator := NewListUsersPaginator(client, nil)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("%v", err)
		}

		users = append(users, output.Users...)
	}

	if len(users) != 45 {
		t.Errorf("ListUsersPaginator(), got = %d, want = %d", len(users), 45)
	}
}

func mustURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()

	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("mustURL() error = %v", err)
	}

	return u
}

func newUsers(t *testing.T, num int) []User {
	t.Helper()

	var users []User
	for i := 0; i < num; i++ {
		users = append(users, User{BaseUser{Id: fmt.Sprintf("%d", i)}})
	}

	return users
}

func Test_newListUsersURL(t *testing.T) {
	type args struct {
		params *ListUsersInput
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			"no-settings",
			args{params: &ListUsersInput{}},
			mustURL(t, fmt.Sprintf("%s/users", EndpointProduction)),
			false,
		},
		{
			"limit",
			args{params: &ListUsersInput{Limit: 25}},
			mustURL(t, fmt.Sprintf("%s/users?limit=25", EndpointProduction)),
			false,
		},
		{
			"skip",
			args{params: &ListUsersInput{Skip: 7}},
			mustURL(t, fmt.Sprintf("%s/users?skip=7", EndpointProduction)),
			false,
		},
		{
			"email",
			args{params: &ListUsersInput{Email: "test@example.com"}},
			mustURL(t, fmt.Sprintf("%s/users?email=%s", EndpointProduction, url.QueryEscape("test@example.com"))),
			false,
		},
		{
			"custom_property",
			args{params: &ListUsersInput{CustomPropertyName: "department=marketing"}},
			mustURL(t, fmt.Sprintf("%s/users?custom_property_name=%s", EndpointProduction, url.QueryEscape("department=marketing"))),
			false,
		},
		{
			"sortby",
			args{params: &ListUsersInput{SortBy: SortPropertyCountry}},
			mustURL(t, fmt.Sprintf("%s/users?sort=%s", EndpointProduction, SortPropertyCountry)),
			false,
		},
		{
			"sortby-descending",
			args{params: &ListUsersInput{SortBy: SortPropertyCountry, SortOrder: SortOrderDescending}},
			mustURL(t, fmt.Sprintf("%s/users?sort=-%s", EndpointProduction, SortPropertyCountry)),
			false,
		},
		{
			"sortby-ascending",
			args{params: &ListUsersInput{SortBy: SortPropertyCountry, SortOrder: SortOrderAscending}},
			mustURL(t, fmt.Sprintf("%s/users?sort=%s", EndpointProduction, SortPropertyCountry)),
			false,
		},
		{
			"include_archived-true",
			args{params: &ListUsersInput{IncludeArchived: true}},
			mustURL(t, fmt.Sprintf("%s/users?include_archived=true", EndpointProduction)),
			false,
		},
		{
			"include_archived-false",
			args{params: &ListUsersInput{IncludeArchived: false}},
			mustURL(t, fmt.Sprintf("%s/users", EndpointProduction)),
			false,
		},
		{
			"show_financial_data-true",
			args{params: &ListUsersInput{ShowFinancialData: true}},
			mustURL(t, fmt.Sprintf("%s/users?show_financial_data=true", EndpointProduction)),
			false,
		},
		{
			"show_financial_data-false",
			args{params: &ListUsersInput{ShowFinancialData: false}},
			mustURL(t, fmt.Sprintf("%s/users", EndpointProduction)),
			false,
		},
		{
			"usermode",
			args{params: &ListUsersInput{UserMode: UserModeBenefactor}},
			mustURL(t, fmt.Sprintf("%s/users?user_mode=%s", EndpointProduction, UserModeBenefactor)),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newListUsersURL(EndpointProduction, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("newListUsersURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("newListUsersURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
