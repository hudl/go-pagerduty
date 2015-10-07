package pagerduty

import (
	"encoding/json"
	"io"
	"time"
)

const (
	WebhookIncidentAcknowledge   = "incident.acknowledge"
	WebhookIncidentAssign        = "incident.assign"
	WebhookIncidentDelegate      = "incident.delegate"
	WebhookIncidentEscalate      = "incident.escalate"
	WebhookIncidentResolve       = "incident.resolve"
	WebhookIncidentTrigger       = "incident.trigger"
	WebhookIncidentUnacknowledge = "incident.unacknowledge"
)

// WebhooksService provides functionality to interact with messages received
// from the PagerDuty Webhooks API.
type WebhooksService struct {
	client *Client
}

type WebhookIncident struct {
	ID                    *string                `json:"id"`
	Number                *int                   `json:"incident_number"`
	CreatedOn             *time.Time             `json:"created_on"`
	Status                *string                `json:"status"`
	URL                   *string                `json:"html_url"`
	Key                   *string                `json:"incident_key"`
	Service               *Service               `json:"service"`
	AssignedToUser        *User                  `json:"assigned_to_user"`
	TriggerSummaryData    map[string]interface{} `json:"trigger_summary_data"`
	TriggerDetailsHTMLURL *string                `json:"trigger_details_html_url"`
	LastStatusChangeOn    *time.Time             `json:"last_status_change_on"`
	LastStatusChangeBy    *string                `json:"last_status_change_by"`
}

type webhookIncidentWrapper struct {
	Incident *WebhookIncident `json:"incident"`
}

type WebhookMessage struct {
	ID        *string                 `json:"id"`
	Type      *string                 `json:"type"`
	CreatedOn *time.Time              `json:"created_on"`
	Data      *webhookIncidentWrapper `json:"data"`
}

type webhookMessageListWrapper struct {
	Messages []WebhookMessage `json:"messages"`
}

func (s *WebhooksService) DecodeMessages(reader io.Reader) ([]WebhookMessage, error) {
	messages := new(webhookMessageListWrapper)
	err := json.NewDecoder(reader).Decode(messages)
	if err != nil {
		return nil, err
	}

	return messages.Messages, err
}
