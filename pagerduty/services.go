package pagerduty

import (
	"fmt"
	"time"
)

const (
	StatusActive      = "active"
	StatusWarning     = "warning"
	StatusCritical    = "critical"
	StatusMaintenance = "maintenance"
	StatusDisabled    = "disabled"

	ServiceTypeCloudkick        = "cloudkick"
	ServiceTypeGenericEmail     = "generic_email"
	ServiceTypeGenericEventsAPI = "generic_events_api"
	ServiceTypeKeynote          = "keynote"
	ServiceTypeNagios           = "nagios"
	ServiceTypePingdom          = "pingdom"
	ServiceTypeServerDensity    = "server_density"
	ServiceTypeSQLMonitor       = "sql_monitor"
)

// ServicesService handles communication with the Services related methods of
// the PagerDuty API.
type ServicesService struct {
	client *Client
}

// Service represents a PagerDuty service.
type Service struct {
	ID                     *string           `json:"id,omitempty"`
	Name                   *string           `json:"name,omitempty"`
	Description            *string           `json:"description,omitempty"`
	URL                    *string           `json:"service_url,omitempty"`
	Key                    *string           `json:"service_key,omitempty"`
	AutoResolveTimeout     *int              `json:"auto_resolve_timeout,omitempty"`
	AcknowledgementTimeout *int              `json:"acknowledgement_timeout,omitempty"`
	CreatedAt              *time.Time        `json:"created_at,omitempty"`
	Status                 *string           `json:"status,omitempty"`
	LastIncidentTimestamp  *time.Time        `json:"last_incident_timestamp,omitempty"`
	EmailIncidentCreation  *string           `json:"email_incident_creation,omitempty"`
	IncidentCounts         *IncidentCounts   `json:"incident_counts,omitempty"`
	EmailFilterMode        *string           `json:"email_filter_mode,omitempty"`
	Type                   *string           `json:"service_type,omitempty"`
	EscalationPolicy       *EscalationPolicy `json:"escalation_policy,omitempty"`
	EmailFilters           []EmailFilter     `json:"email_filters,omitempty"`
	SeverityFilter         *string           `json:"severity_filter,omitempty"`
}

type IncidentCounts struct {
	Triggered    *int `json:"triggered,omitempty"`
	Acknowledged *int `json:"acknowledged,omitempty"`
	Resolved     *int `json:"resolved,omitempty"`
	Total        *int `json:"total,omitempty"`
}

type EmailFilter map[string]interface{}

type ServiceListOptions struct {
	// A comma-separated list of team IDs, specifying teams whose maintenance
	// windows will be returned.
	Teams string `url:"teams,omitempty"`

	// A comma-separated list of extra fields to include in the response. Valid
	// fields include 'escalation_policy', 'email_filters', and 'teams'.
	Include []string `url:"include,omitempty"`

	// Time zone in which the dates in the result will be redered. Defaults to
	// account default time zone.
	TimeZone *TimeZone `url:"time_zone,omitempty"`

	// Filters the result, showing only services whose 'name' or 'service_key'
	// matches the query.
	Query string `url:"query,omitempty"`

	// Specifies the field to sort the response on, defaults to 'name'. Valid
	// fields are 'name' and 'id'.
	SortBy string `url:"sort_by,omitempty"`

	ListOptions
}

type serviceListWrapper struct {
	Services []Service `json:"services"`
}

// List services filtered by provided options.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/services/list
func (s *ServicesService) List(opts *ServiceListOptions) ([]Service, *Response, error) {
	uri, err := addOptions("services", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	services := new(serviceListWrapper)
	resp, err := s.client.Do(req, services)
	if err != nil {
		return nil, resp, err
	}

	return services.Services, resp, err
}

type GetServiceOptions struct {
	// Include extra information in the response. Possible values are
	// 'escalation_policy' and 'email_filters'.
	Include []string `url:"include,omitempty"`
}

// Get fetches a service by id and filtered by provided options.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/services/show
func (s *ServicesService) Get(id string, opts *GetServiceOptions) (*Service, *Response, error) {
	path := fmt.Sprintf("services/%s", id)
	uri, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	service := new(Service)
	resp, err := s.client.Do(req, service)
	if err != nil {
		return nil, resp, err
	}

	return service, resp, err
}
