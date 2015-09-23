package pagerduty_test

import (
	. "github.com/hudl/go-pagerduty/pagerduty"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"net/http"
	"net/url"
	"testing"
)

const (
	// The default test base URL for the PagerDuty API.
	defaultBaseURL = "https://" + subdomain + ".pagerduty.com/api/v1/"

	// The test subdomain.
	subdomain = "subdomain"

	// The test api key.
	apiKey = "super-secret-key"
)

// TestEnvironment is a complete testing environment for mocking the PagerDuty
// API server.
type TestEnvironment struct {
	// Server is a ghttp server used to provide mock API responses.
	Server *ghttp.Server

	// Client is the PagerDuty client being tested
	Client *Client
}

// NewTestEnvironment creates and configures a new test environment with a
// ghttp test server along with a PagerDuty client to talk to the server. Tests
// should register handlers on the server which provide mock responses for the
// API method being tested.
func NewTestEnvironment() *TestEnvironment {
	server := ghttp.NewServer()

	// PagerDuty client configured to use the test server
	client := NewClient(nil, subdomain, apiKey)
	url, _ := url.Parse(server.URL())
	client.BaseURL = url

	return &TestEnvironment{
		Server: server,
		Client: client,
	}
}

// verifyHeaderHandler is an http.HandlerFunc that checks for the proper
// headers in a request.
var verifyHeaderHandler = ghttp.CombineHandlers(
	ghttp.VerifyHeader(http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Token token=" + apiKey},
	}),
	ghttp.VerifyContentType("application/json"),
)

func TestPagerDuty(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PagerDuty Suite")
}