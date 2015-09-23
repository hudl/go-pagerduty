package pagerduty_test

import (
	. "github.com/hudl/go-pagerduty/pagerduty"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"encoding/json"
	"net/http"
	"regexp"
	// "time"
)

const (
	incidentListJSON = `{ "incidents": [` + incidentJSON + `]}`
	incidentJSON     = `{
		"incident_number": 1 ,
		"created_on": null,
		"status": "status",
		"html_url": "html_url",
		"incident_key": null,
		"pending_actions": [{
			"type": "type",
			"at": null
		}],
		"service": {
			"id": "id",
			"name": "name",
			"description": "description",
			"html_url": "html_url"
		},
		"assigned_to_user": {
			"id": "id",
			"name": "name",
			"email": "email",
			"html_url": "html_url"
		},
		"assigned_to": [{
			"at": null,
			"object": {
				"type": "type"
			}
		}],
		"trigger_summary_data": {
			"subject": "subject"
		},
		"trigger_details_html_url": "trigger_details_html_url",
		"last_status_change_on": null,
		"last_status_change_by": null,
		"urgency": "urgency"
	}`
)

var _ = Describe("Incidents", func() {
	var (
		env              *TestEnvironment
		expectedIncident Incident
	)

	json.Unmarshal([]byte(incidentJSON), &expectedIncident)

	BeforeEach(func() { env = NewTestEnvironment() })
	AfterEach(func() { env.Server.Close() })

	Describe("List", func() {
		Context("with a successful, non-empty response", func() {
			var (
				incidents []Incident
				resp      *Response
				err       error
			)

			BeforeEach(func() {
				env.Server.RouteToHandler(GET, "/incidents", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, incidentListJSON),
				))

				incidents, resp, err = env.Client.Incidents.List(nil)
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return a non-empty response", func() {
				Expect(resp).NotTo(BeNil())
			})

			It("should return the expected incidents", func() {
				Expect(incidents).NotTo(BeNil())
				Expect(incidents).NotTo(BeEmpty())
				Expect(incidents[0]).To(Equal(expectedIncident))
			})
		})
	})

	Describe("Get", func() {
		Context("with a successful, non-empty response", func() {
			var (
				incident *Incident
				resp     *Response
				err      error
			)

			BeforeEach(func() {
				path, _ := regexp.Compile("/incidents/\\w+")
				env.Server.RouteToHandler(GET, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, incidentJSON),
				))

				incident, resp, err = env.Client.Incidents.Get("id")
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return a non-empty response", func() {
				Expect(resp).NotTo(BeNil())
			})

			It("should return the expected incident", func() {
				Expect(incident).NotTo(BeNil())
				Expect(incident).To(Equal(&expectedIncident))
			})
		})
	})
})
