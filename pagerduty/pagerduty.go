// This package was modeled after Google's awesome go-github library
//   https://github.com/google/go-github
//
// This file contains a few helper functions from the project.

package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL   = "https://%s.pagerduty.com/api/v1/"
	defaultEventsURL = "https://events.pagerduty.com/"

	headerAuthorization = "Authorization"
	headerAccept        = "Accept"
	headerContentType   = "Content-Type"

	authorizationToken = "Token token=%s"
	acceptType         = "application/json"
	contentType        = "application/json"

	// http verb constants
	DELETE = "DELETE"
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
)

// A Client manages communication with the PagerDuty API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for the PagerDuty API.
	BaseURL *url.URL

	// Base URL for the PagerDuty Events API.
	EventsURL *url.URL

	// Client subdomain for the PagerDuty API.
	subdomain string

	// PagerDuty API key.
	APIKey string

	// Services used for talking to different parts of the PagerDuty API.
	Alerts             *AlertsService
	EscalationPolicies *EscalationPoliciesService
	Events             *EventsService
	Incidents          *IncidentsService
	Schedules          *SchedulesService
	Services           *ServicesService
	Teams              *TeamsService
	Users              *UsersService
	Webhooks           *WebhooksService
}

// NewClient returns a new PagerDuty API client. If httpClient is nil,
// http.DefaultClient will be used.
func NewClient(httpClient *http.Client, subdomain string, apiKey ...string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	renderedURL := fmt.Sprintf(defaultBaseURL, subdomain)
	baseURL, _ := url.Parse(renderedURL)
	eventsURL, _ := url.Parse(defaultEventsURL)
	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		EventsURL: eventsURL,
		subdomain: subdomain,
	}

	if len(apiKey) == 1 {
		c.APIKey = apiKey[0]
	}

	// configure services
	c.Alerts = &AlertsService{client: c}
	c.EscalationPolicies = &EscalationPoliciesService{client: c}
	c.Events = &EventsService{client: c}
	c.Incidents = &IncidentsService{client: c}
	c.Schedules = &SchedulesService{client: c}
	c.Services = &ServicesService{client: c}
	c.Teams = &TeamsService{client: c}
	c.Users = &UsersService{client: c}
	c.Webhooks = &WebhooksService{client: c}

	return c
}

type ListOptions struct {
	// The offset of the first record returned. Default is 0.
	Offset int `url:"offset,omitempty"`

	// The number of records returned. Default (and max limit) is 100 for most APIs.
	Limit int `url:"limit,omitempty"`
}

// addOptions adds the parameters in opt as URL query parameters to s.
// opt must be a struct whose fields may contain "url" tags.
//
// This function is from the google/go-github project
// https://github.com/google/go-github/blob/7277108aa3e8823e0e028f6c74aea2f4ce4a1b5a/github/github.go#L102-L122
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// newRequest is a helper function to generate an http.Request and automagically
// json encode a body and resolve the given path.
func newRequest(baseURL *url.URL, method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	uri := baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, uri.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add(headerAccept, acceptType)
	req.Header.Add(headerContentType, contentType)

	return req, nil
}

// NewRequest creates an API request. A relative URL can be provided in path,
// in which case it is resolved relative to the BaseURL of the client.
// Relative URLs should always be specified without the preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	req, err := newRequest(c.BaseURL, method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add(headerAuthorization, fmt.Sprintf(authorizationToken, c.APIKey))

	return req, nil
}

// Response is a PagerDuty API response. It wraps the standard http.Response
// returned from PagerDuty and provides convinient access to pagination
// response fields.
type Response struct {
	*http.Response

	// Pagination response fields. Any or all of these may be set to the zero
	// zero value for responses that are not part of a paginated set.

	Offset int `json:"offset,omitempty"` // The offset used in the execution of the query
	Limit  int `json:"limit,omitempty"`  // The limit used in the execution of the query
	Total  int `json:"total,omitemtpy"`  // The total number of records available
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	r := &Response{Response: resp}

	errResp := CheckResponse(resp)

	// save the response body so it can be unmarshalled multiple times
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	// body is not empty
	if len(body) > 0 {
		// unmarshal pagination response fields
		err = json.Unmarshal(body, r)
		if err != nil {
			return r, err
		}

		if v != nil {
			err = json.Unmarshal(body, v)
			if err != nil {
				return r, err
			}
		}
	}

	return r, errResp
}

// Delete is a convenience function to create and execute a DELETE request.
func (c *Client) Delete(path string) (*Response, error) {
	req, err := c.NewRequest(DELETE, path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// Get is a convenience function to create and execute a GET request.
func (c *Client) Get(path string, v interface{}) (*Response, error) {
	req, err := c.NewRequest(GET, path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, v)
}

// Post is a convenience function to create and execute a POST request.
func (c *Client) Post(path string, body, v interface{}) (*Response, error) {
	req, err := c.NewRequest(POST, path, body)
	if err != nil {
		return nil, err
	}

	return c.Do(req, v)
}

// Put is a convenience function to create and execute a PUT request.
func (c *Client) Put(path string, body, v interface{}) (*Response, error) {
	req, err := c.NewRequest(PUT, path, body)
	if err != nil {
		return nil, err
	}

	return c.Do(req, v)
}

// An ErrorResponse reports an error caused by an API request.
type ErrorResponse struct {
	Response *http.Response
	Code     int      `json:"code,omitempty"`
	Message  string   `json:"message,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e.Response != nil && (e.Response.StatusCode < 200 || e.Response.StatusCode > 299) {
		return fmt.Sprintf("pagerduty: %v %q: %d",
			e.Response.Request.Method,
			e.Response.Request.URL,
			e.Response.StatusCode)
	}

	return fmt.Sprintf("pagerduty: api error %d %q", e.Code, e.Message)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range.
func CheckResponse(r *http.Response) error {
	c := r.StatusCode
	if 200 <= c && c <= 299 {
		return nil
	}

	er := &ErrorResponse{Response: r}
	body, err := ioutil.ReadAll(r.Body)
	if err == nil && body != nil {
		json.Unmarshal(body, er)
	}

	return er
}

// These helpers are from google's go-github
// https://github.com/google/go-github/blob/7277108aa3e8823e0e028f6c74aea2f4ce4a1b5a/github/github.go#L565-L588

// Bool is a helper function that allocates a new bool value to store v and
// returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// Int is a helper function that allocates a new int value to store v and
// returns a pointer to it.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// String is a helper function that allocates a new string value to store v and
// returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Time is a helper function that allocates a new time.Time value to store v
// and returns a pointer to it.
func Time(v time.Time) *time.Time {
	p := new(time.Time)
	*p = v
	return p
}
