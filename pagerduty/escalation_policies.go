package pagerduty

import (
	"fmt"
)

// EscalationPoliciesService handles communication with the Escalation Policies
// related methods of the PagerDuty API.
type EscalationPoliciesService struct {
	client *Client
}

// EscalationPolicy represents a PagerDuty escalation policy.
type EscalationPolicy struct {
	ID              *string          `json:"id,omitempty"`
	Name            *string          `json:"name,omitempty"`
	NumLoops        *int             `json:"num_loops,omitempty"`
	EscalationRules []EscalationRule `json:"escalation_rules,omitempty"`
	Services        []Service        `json:"services,omitempty"`
}

type EscalationRule struct {
	ID                       *string  `json:"id,omitempty"`
	EscalationDelayInMinutes *int     `json:"escalation_delay_in_minutes,omitempty"`
	Targets                  []Target `json:"targets,omitempty"`
}

type Target struct {
	ID       *string   `json:"id,omitempty"`
	Type     *string   `json:"type,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Email    *string   `json:"email,omitempty"`
	TimeZone *TimeZone `json:"time_zone,omitempty"`
	Color    *string   `json:"color,omitempty"`
}

type EscalationPolicyListOptions struct {
	// Filters the results, showing only the escaltion policies whose names
	// match the query.
	Query string `url:"query,omitempty"`

	// A comma-separated lists of team IDs, speicifying teams whose maintenance
	// windows will be returned.
	Teams string `url:"teams,omitempty"`

	// Include extra information in the response. Possible values are 'teams'.
	Include []string `url:"teams,omitempty"`

	ListOptions
}

type escalationPolicyListWrapper struct {
	EscalationPolicies []EscalationPolicy `json:"escalation_policies"`
}

// List escalation policies filtered by provided options.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/escalation_policies/list
func (s *EscalationPoliciesService) List(opts *EscalationPolicyListOptions) ([]EscalationPolicy, *Response, error) {
	uri, err := addOptions("escalation_policies", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	policies := new(escalationPolicyListWrapper)
	resp, err := s.client.Do(req, policies)
	if err != nil {
		return nil, resp, err
	}

	return policies.EscalationPolicies, resp, err
}

type escalationPolicyWrapper struct {
	EscalationPolicy *EscalationPolicy `json:"escalation_policy"`
}

// Get fetches an escalation policy by id.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/escalation_policies/show
func (s *EscalationPoliciesService) Get(id string) (*EscalationPolicy, *Response, error) {
	uri := fmt.Sprintf("escalation_policies/%s", id)

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	policy := new(escalationPolicyWrapper)
	resp, err := s.client.Do(req, policy)
	if err != nil {
		return nil, resp, err
	}

	return policy.EscalationPolicy, resp, err
}
