package pagerduty

import (
	"fmt"
	"time"
)

// SchedulesService handles communication with the Schedules related methods
// of the PagerDuty API.
type SchedulesService struct {
	client *Client
}

// Schedule represents a PagerDuty schedule.
type Schedule struct {
	ID                   *string            `json:"id,omitempty"`
	Name                 *string            `json:"name,omitempty"`
	TimeZone             *TimeZone          `json:"time_zone,omitempty"`
	Today                *Date              `json:"today,omitempty"`
	EscalationPolicies   []EscalationPolicy `json:"escalation_policies,omitempty"`
	ScheduleLayers       []ScheduleLayer    `json:"schedule_layers,omitempty"`
	OverridesSubschedule *ScheduleLayer     `json:"overrides_subschedule,omitempty"`
	FinalSchedule        *ScheduleLayer     `json:"final_schedule,omitempty"`
}

// ScheduleLayer represents one of potentially many layers for a PagerDuty
// schedule.
type ScheduleLayer struct {
	ID                         *string    `json:"id,omitempty"`
	Name                       *string    `json:"name,omitempty"`
	Priority                   *int       `json:"priority,omitempty"`
	Start                      *Date      `json:"start,omitempty"`
	End                        *Date      `json:"end,omitempty"`
	Users                      []User     `json:"users,omitempty"`
	RenderedScheduleEntries    []string   `json:"rendered_schedule_entries,omitempty"`
	RestrictionType            *string    `json:"restriction_type,omitempty"`
	Restrictions               []string   `json:"restrictions,omitempty"`
	RenderedCoveragePercentage *int       `json:"rendered_coverage_percentage,omitempty"`
	RotationTurnLengthSeconds  *int       `json:"rotation_turn_length_seconds,omitempty"`
	RotationVirtualStart       *time.Time `json:"rotation_virtual_start,omitempty"`
}

type ScheduleListOptions struct {
	// Filters the result, showing only the schedules whose name matches the
	// query.
	Query string `url:"query,omitempty"`

	// The user id of the user making the request. This will be used to
	// generate the calendar private urls. This is only needed if you are
	// using token based authentication.
	RequesterID string `url:"requester_id,omitempty"`

	ListOptions
}

type scheduleListWrapper struct {
	Schedules []Schedule `json:"schedules"`
}

// List schedules filtered by provided options
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/schedules/list
func (s *SchedulesService) List(opts *ScheduleListOptions) ([]Schedule, *Response, error) {
	uri, err := addOptions("schedules", opts)
	if err != nil {
		return nil, nil, err
	}

	schedules := new(scheduleListWrapper)
	resp, err := s.client.Get(uri, schedules)
	if err != nil {
		return nil, resp, err
	}

	return schedules.Schedules, resp, err
}

type scheduleWrapper struct {
	Schedule *Schedule `json:"schedule"`
}

// Get fetches a schedule by id.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/schedules/show
func (s *SchedulesService) Get(id string) (*Schedule, *Response, error) {
	uri := fmt.Sprintf("schedules/%s", id)

	schedule := new(scheduleWrapper)
	resp, err := s.client.Get(uri, schedule)
	if err != nil {
		return nil, resp, err
	}

	return schedule.Schedule, resp, err
}
