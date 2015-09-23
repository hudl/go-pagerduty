package pagerduty

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	StatusTriggered    = "triggered"
	StatusAcknowledged = "acknoledged"
	StatusResolved     = "resolved"

	ObjectTypeUser = "user"
	ObjectTypeAPI  = "api"
)

// IncidentsService handles communication with the Incidents related methods of
// the PagerDuty APIl.
type IncidentsService struct {
	client *Client
}

// Incident represents a PagerDuty incident.
type Incident struct {
	Number             *int              `json:"incident_number,omitempty"`
	Status             *string           `json:"status,omitempty"`
	Urgency            *string           `json:"urgency,omitempty"`
	PendingActions     []PendingAction   `json:"pending_actions,omitempty"`
	CreatedOn          *time.Time        `json:"created_on,omitempty"`
	URL                *string           `json:"html_url,omitempty"`
	Key                *string           `json:"incident_key,omitempty"`
	Service            *Service          `json:"service,omitempty"`
	EscalationPolicy   *EscalationPolicy `json:"escalation_policy,omitempty"`
	Teams              []Team            `json:"teams,omitempty"`
	AssignedTo         []ObjectAt        `json:"assigned_to,omitempty"`
	Acknowledgers      []ObjectAt        `json:"acknowledgers,omitempty"`
	LastStatusChangeBy *User             `json:"last_status_change_by,omitempty"`
	LastStatusChangeOn *time.Time        `json:"last_status_chage_on,omitempty"`
	TriggerSummary     *TriggerSummary   `json:"trigger_summary_data,omitempty"`
	TriggerDetailsURL  *string           `json:"trigger_details_html_url,omitempty"`
}

type PendingAction struct {
	Type *string    `json:"type"`
	At   *time.Time `json:"at"`
}

type ObjectAt struct {
	At   *time.Time `json:"at,omitempty"`
	Type *string    `json:"-"`

	// original message object
	Object map[string]interface{} `json:"object,omitempty"`
}

func (a *ObjectAt) UnmarshalJSON(data []byte) error {
	type alias ObjectAt
	temp := new(alias)
	err := json.Unmarshal(data, temp)
	if err != nil {
		return err
	}

	*a = ObjectAt(*temp)
	a.Type = String(a.Object["type"].(string))
	return nil
}

type TriggerSummary map[string]interface{}

type IncidentListOptions struct {
	// The start of the date range you want to search.
	Since *time.Time `url:"since,omitempty"`

	// The end of the date range you want to search.
	Until *time.Time `url:"until,omitemtpy"`

	// When set to 'all', the 'since' and 'until' parameters and defaults are
	// ignored. Unse this to get all incidents since the account was created.
	DateRange string `url:"date_range,omitempty"`

	// Used to restrict the properties of each incident returned to a set of
	// pre-defined fields. If ommited, returned incidents have the majority
	// of fields present.
	Fields string `url:"fields,omitempty"`

	// Returns only the incidents in the passed status(es). Valid status
	// options are 'triggered', 'acknowledged' and 'resolved'.
	Status string `url:"status,omitemtpy"`

	// Returns only the incidents with the passes de-duplication key.
	IncidentKey string `url:"incidient_key,omitempty"`

	// Returns only the incidents associated with the passed service(s).
	// Expects one or more service IDs separated by commas.
	Service string `url:"service,omitempty"`

	// A comma-separated list of team IDs, specifying teams whose maintenance
	// windows will be returned.
	Teams string `url:"teams,omitempty"`

	// Returns only the incidents currently assigned to the passed user(s).
	// This expects one or more user IDs separated by commas.
	//
	// NOTE: When using the 'assigned_to_user' filter, you will only receive
	// incidents with statuses of 'triggered' or 'acknowledged', because
	// 'resolved' incidents are not assigned to any user.
	AssignedToUser string `url:"assigned_to_user,omitempty"`

	// A comma-separated list of urgencies to filter the incidents list. Defaults
	// to 'high,low'.
	Urgency string `url:"urgency,omitempty"`

	// Time zones in which dates in the result will be rendered. Defaults to
	// 'UTC'.
	TimeZone *TimeZone `url:"time_zone,omitempty"`

	// A comma-separated list of fields in which to sort the results, as well as
	// the direction (ascending/descending).
	//
	// Valid fields are 'incident_number', 'created_on', 'resolved_on' and
	// 'urgency'.
	SortBy string `url:"sort_by,omitempty"`

	ListOptions
}

type incidentListWrapper struct {
	Incidents []Incident `json:"incidents"`
}

// List incidents filtered by provided options.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/incidents/list
func (s *IncidentsService) List(opts *IncidentListOptions) ([]Incident, *Response, error) {
	uri, err := addOptions("incidents", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	incidents := new(incidentListWrapper)
	resp, err := s.client.Do(req, incidents)
	if err != nil {
		return nil, resp, err
	}

	return incidents.Incidents, resp, err
}

// Get fetches an incident by id.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/incidents/show
func (s *IncidentsService) Get(id string) (*Incident, *Response, error) {
	uri := fmt.Sprintf("incidents/%s", id)

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	incident := new(Incident)
	resp, err := s.client.Do(req, incident)
	if err != nil {
		return nil, resp, err
	}

	return incident, resp, err
}
