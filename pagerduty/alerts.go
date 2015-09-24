package pagerduty

import (
	"time"
)

const (
	AlertTypeEmail = "Email"
	AlertTypePhone = "Phone"
	AlertTypePush  = "Push"
	AlertTypeSMS   = "SMS"
)

// AlertsService handles communication with the Alerts related methods of the
// PagerDuty API.
type AlertsService struct {
	client *Client
}

// Alert represents a PagerDuty alert.
type Alert struct {
	ID        *string    `json:"id"`
	Type      *string    `json:"type"`
	StartedAt *time.Time `json:"started_at"`
	User      *User      `json:"user"`
	Address   *string    `json:"address"`
}

type AlertListOptions struct {
	// The start of the date range you want to search.
	Since time.Time `url:"since,omitempty"`

	// The end of the date range you want to search.
	Until time.Time `url:"until,omitempty"`

	// Returns only alerts of the said types. Can be one of 'SMS', 'Email',
	// 'Pone' or 'Push'.
	FilterType string `url:"filter[type],omitempty"`

	// Time zone in which dates in the result will be rendered. Defaults to
	// account time zone.
	TimeZone TimeZone `url:"time_zone,omitempty"`

	ListOptions
}

type alertListWrapper struct {
	Alerts []Alert `json:"alerts"`
}

// List alerts filtered by provided options.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/alerts/list
func (s *AlertsService) List(opts *AlertListOptions) ([]Alert, *Response, error) {
	uri, err := addOptions("alerts", opts)
	if err != nil {
		return nil, nil, err
	}

	alerts := new(alertListWrapper)
	resp, err := s.client.Get(uri, alerts)
	if err != nil {
		return nil, resp, err
	}

	return alerts.Alerts, resp, err
}
