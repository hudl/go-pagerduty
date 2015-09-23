package pagerduty

import (
	"fmt"
)

// TeamsService handles communication with the Teams related methods of the
// PagerDuty API.
type TeamsService struct {
	client *Client
}

// Team represents a PagerDuty team.
type Team struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type TeamListOptions struct {
	// Filters the results, showing only teams whose names match the query.
	Query string `url:"query,omitempty"`

	ListOptions
}

type teamListWrapper struct {
	Teams []Team `json:"teams"`
}

// List teams filtered by provided options.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/teams/list
func (s *TeamsService) List(opts *TeamListOptions) ([]Team, *Response, error) {
	uri, err := addOptions("teams", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	teams := new(teamListWrapper)
	resp, err := s.client.Do(req, teams)
	if err != nil {
		return nil, resp, err
	}

	return teams.Teams, resp, err
}

type teamWrapper struct {
	Team *Team `json:"team"`
}

// Get fetches a team by id.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/teams/show
func (s *TeamsService) Get(id string) (*Team, *Response, error) {
	uri := fmt.Sprintf("teams/%s", id)

	req, err := s.client.NewRequest(GET, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	team := new(teamWrapper)
	resp, err := s.client.Do(req, team)
	if err != nil {
		return nil, resp, err
	}

	return team.Team, resp, err
}

// Create a team.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/teams/create
func (s *TeamsService) Create(team *Team) (*Team, *Response, error) {
	uri := "teams"

	req, err := s.client.NewRequest(POST, uri, team)
	if err != nil {
		return nil, nil, err
	}

	t := new(Team)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

// Edit a team.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/teams/update
func (s *TeamsService) Edit(team *Team) (*Team, *Response, error) {
	uri := fmt.Sprintf("teams/%s", team.ID)

	req, err := s.client.NewRequest(PUT, uri, team)
	if err != nil {
		return nil, nil, err
	}

	t := new(Team)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

// Delete a team.
//
// PagerDuty API docs: https://developer.pagerduty.com/documentation/rest/teams/delete
func (s *TeamsService) Delete(id string) (*Response, error) {
	uri := fmt.Sprintf("teams/%s", id)

	req, err := s.client.NewRequest(DELETE, uri, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
