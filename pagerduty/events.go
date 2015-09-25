package pagerduty

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	eventsAPIPath = "generic/2010-04-15/create_event.json"

	ContextTypeLink  = "link"
	ContextTypeImage = "image"

	EventTypeAcknowledge = "acknowledge"
	EventTypeResolve     = "resolve"
	EventTypeTrigger     = "trigger"
)

// NewEventRequest creates an API request for the PagerDuty Events API. It has
// the same requirements and behaviors as Client.NewRequest.
func (c *Client) NewEventRequest(method, path string, body interface{}) (*http.Request, error) {
	return newRequest(c.EventsURL, method, path, body)
}

type EventResponse struct {
	Response *http.Response

	Status      string   `json:"status,omitempty"`
	Message     string   `json:"message,omitempty"`
	IncidentKey string   `json:"incident_key,omitempty"`
	Errors      []string `json:"errors,omitempty"`
}

// DoEventRequest sends a request to the PagerDuty Events API and returns the
// response.
func (c *Client) DoEventRequest(req *http.Request) (*EventResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	eventResp := new(EventResponse)
	err = json.NewDecoder(resp.Body).Decode(eventResp)
	if err != nil {
		return nil, err
	}

	eventResp.Response = resp
	return eventResp, nil
}

// EventsService handles communication with the PagerDuty Events API.
//
// https://developer.pagerduty.com/documentation/integration/events
type EventsService struct {
	client *Client
}

type Event struct {
	Type        *string                `json:"event_type,omitempty"`
	ServiceKey  *string                `json:"service_key,omitempty"`
	Description *string                `json:"description,omitempty"`
	IncidentKey *string                `json:"incident_key,omitempty"`
	Client      *string                `json:"client,omitempty"`
	ClientURL   *string                `json:"client_url,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`

	// NOTE: Contexts should contain only LinkContext or ImageContext structs.
	Contexts []interface{} `json:"contexts,omitempty"`
}

type LinkContext struct {
	// The type of context being attached to the incident. Possible values are
	// 'link' and 'image'.
	Type string `json:"type,omitempty"`

	// The link to either the incident being attached or image.
	Href string `json:"href,omitempty"`

	// Options information pertaining to the incident.
	Text string `json:"test,omitempty"`
}

type ImageContext struct {
	// The type of context being attached to the incident. Possible values are
	// 'link' and 'image'.
	Type string `json:"type,omitempty"`

	// The source of the image being attached to the incident. This must be
	// served via HTTPS.
	Source string `json:"src,omitempty"`

	// Optional link for the image.
	Href string `json:"href,omitempty"`

	// Optional alternative text for the image.
	Alt string `json:"alt,omitempty"`
}

// helper function to post an event and unmarshal the event response.
func (s *EventsService) postEvent(event *Event, eventType string) (*EventResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("pagerduty: event cannot be nil")
	}

	event.Type = String(eventType)
	req, err := s.client.NewEventRequest(POST, eventsAPIPath, event)
	if err != nil {
		return nil, err
	}

	return s.client.DoEventRequest(req)
}

// Acknowledge an event (incident).
//
// https://developer.pagerduty.com/documentation/integration/events/acknowledge
func (s *EventsService) Acknowledge(event *Event) (*EventResponse, error) {
	return s.postEvent(event, EventTypeAcknowledge)
}

// Resolve an event (incident).
//
// https://developer.pagerduty.com/documentation/integration/events/resolve
func (s *EventsService) Resolve(event *Event) (*EventResponse, error) {
	return s.postEvent(event, EventTypeResolve)
}

// Trigger an event (incident).
//
// https://developer.pagerduty.com/documentation/integration/events/trigger
func (s *EventsService) Trigger(event *Event) (*EventResponse, error) {
	return s.postEvent(event, EventTypeTrigger)
}
