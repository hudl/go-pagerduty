package pagerduty_test

import (
	. "github.com/hudl/go-pagerduty/pagerduty"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"io"
	"strings"
)

const (
	webhookJSON = `{
		"messages": [{
			"id": "bb8b8fe0-e8d5-11e2-9c1e-22000afd16cf",
			"created_on": "2013-07-09T20:25:44Z",
			"type": "incident.trigger",
			"data": {
				"incident": {
					"id": "PIJ90N7",
					"incident_number": 1,
					"created_on": "2013-07-09T20:25:44Z",
					"status": "triggered",
					"html_url": "https://acme.pagerduty.com/incidents/PIJ90N7",
					"incident_key": "null",
					"service": {
						"id": "PBAZLIU",
						"name": "service",
						"html_url": "https://acme.pagerduty.com/services/PBAZLIU"
					},
					"assigned_to_user": {
						"id": "PPI9KUT",
						"name": "Alan Kay",
						"email": "alan@pagerduty.com",
						"html_url": "https://acme.pagerduty.com/users/PPI9KUT"
					},
					"resolved_by_user": {
						"id": "PPI9KUT",
						"name": "Alan Kay",
						"email": "alan@pagerduty.com",
						"html_url": "https://acme.pagerduty.com/users/PPI9KUT"
					},
					"trigger_summary_data": {
						"subject": "45645"
					},
					"trigger_details_html_url": "https://acme.pagerduty.com/incidents/PIJ90N7/log_entries/PIJ90N7",
					"last_status_change_on": "2013-07-09T20:25:44Z",
					"last_status_change_by": {
						"id": "PPI9KUT",
						"name": "Alan Kay",
						"email": "alan@pagerduty.com",
						"html_url": "https://acme.pagerduty.com/users/PPI9KUT"
					}
				}
			}
		}]
	}`
)

type webhookMessageListWrapper struct {
	Messages []WebhookMessage `json:"messages"`
}

var _ = Describe("Webhooks", func() {
	var expectedMessages webhookMessageListWrapper
	json.Unmarshal([]byte(webhookJSON), &expectedMessages)

	Describe("DecodeMessages", func() {
		Context("with a valid, non-empty webhook reader", func() {
			var (
				env      *TestEnvironment
				reader   io.Reader
				messages []WebhookMessage
				err      error
			)

			BeforeEach(func() {
				env = NewTestEnvironment()
				reader = strings.NewReader(webhookJSON)
				messages, err = env.Client.Webhooks.DecodeMessages(reader)
			})

			AfterEach(func() { env.Server.Close() })

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the expected messages", func() {
				Expect(messages).NotTo(BeEmpty())
				Expect(messages).To(Equal(expectedMessages.Messages))
			})
		})
	})
})
