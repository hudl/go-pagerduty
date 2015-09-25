package pagerduty

import (
	"encoding/json"
	"net/http"
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

type WebhookMessage struct {
	ID        *int       `json:"id,omitempty"`
	Type      *string    `json:"type,omitempty"`
	CreatedOn *time.Time `json:"created_on,omitempty"`
	Data      []Incident `json:"data,omitempty"`
}

type webhookMessageListWrapper struct {
	Messages []WebhookMessage `json:"messages"`
}

func (s *WebhooksService) DecodeMessages(resp *http.Response) ([]WebhookMessage, error) {
	messages := new(webhookMessageListWrapper)
	err := json.NewDecoder(resp.Body).Decode(messages)
	if err != nil {
		return nil, err
	}

	return messages.Messages, err
}
