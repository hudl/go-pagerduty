package pagerduty

import (
	"fmt"
)

const (
	IncludeContactMethod     = "contact_method"
	IncludeNotificationRules = "notification_rules"
)

// UsersService handles communication with the Users related methods of the
// PagerDuty API.
type UsersService struct {
	client *Client
}

// User represents a PagerDuty user.
type User struct {
	ID              *string   `json:"id,omitempty"`
	Name            *string   `json:"name,omitempty"`
	Email           *string   `json:"email,omitempty"`
	TimeZone        *TimeZone `json:"time_zone,omitempty"`
	Color           *string   `json:"color,omitempty"`
	Role            *string   `json:"role,omitempty"`
	URL             *string   `json:"user_url,omitempty"`
	AvatarURL       *string   `json:"avatar_url,omitempty"`
	InvitationSent  *bool     `json:"invitation_sent,omitempty"`
	MarketingOptOut *bool     `json:"marketing_opt_out,omitempt"`
	JobTitle        *string   `json:"job_title,omitempty"`
}

type UserListOptions struct {
	// Filters the result, showing only the users whose names or email addresses
	// match the query
	Query string `url:"query,omitemtpy"`

	// Array of additional details to include. This API accepts `contact_method`
	// and `notification_rules`.
	Include []string `url:"include,omitempty"`

	ListOptions
}

type usersListWrapper struct {
	Users []User `json:"users"`
}

// List user filtered by provided options.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/users/list
func (s *UsersService) List(opts *UserListOptions) ([]User, *Response, error) {
	uri, err := addOptions("users", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	users := new(usersListWrapper)
	resp, err := s.client.Do(req, users)
	if err != nil {
		return nil, resp, err
	}

	return users.Users, resp, err
}

type GetUserOptions struct {
	// Array of additional details to include. This API accepts `contact_method`
	// and `notification_rules`.
	Include []string `url:"include,omitempty"`
}

type userWrapper struct {
	User *User `json:"user"`
}

// Get fetches a user by id and filtered by provided options.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/users/show
func (s *UsersService) Get(id string, opts *GetUserOptions) (*User, *Response, error) {
	path := fmt.Sprintf("users/%s", id)
	uri, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(userWrapper)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user.User, resp, err
}
