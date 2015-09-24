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
		var (
			incidents []Incident
			resp      *Response
			err       error
		)

		Context("with a successful, non-empty response", func() {
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

			It("should return a response with the correct status code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("should return the expected incidents", func() {
				Expect(incidents).NotTo(BeNil())
				Expect(incidents).NotTo(BeEmpty())
				Expect(incidents[0]).To(Equal(expectedIncident))
			})
		})
	})

	Describe("Get", func() {
		var (
			incident *Incident
			resp     *Response
			err      error
		)

		Context("with a successful, non-empty response", func() {
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

			It("should return a response with the correct status code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("should return the expected incident", func() {
				Expect(incident).NotTo(BeNil())
				Expect(incident).To(Equal(&expectedIncident))
			})
		})
	})

	Describe("Count", func() {
		var (
			count int
			resp  *Response
			err   error
		)

		Context("with a successful, non-empty response", func() {
			BeforeEach(func() {
				env.Server.RouteToHandler(GET, "/incidents/count", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, `{"total": 1}`),
				))

				count, resp, err = env.Client.Incidents.Count(nil)
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

			It("should return a response with the correct status code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("should return the expected incident count", func() {
				Expect(count).To(Equal(1))
			})
		})
	})

	Describe("Edit", func() {
		var (
			incidents []Incident
			resp      *Response
			err       error
		)

		Context("with a successful, non-empty response", func() {
			BeforeEach(func() {
				env.Server.RouteToHandler(PUT, "/incidents", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, incidentListJSON),
				))

				incidents, resp, err = env.Client.Incidents.Edit(nil)
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

			It("should return a response with the correct status code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("should return the expected incident count", func() {
				Expect(incidents).NotTo(BeNil())
				Expect(incidents).NotTo(BeEmpty())
				Expect(incidents[0]).To(Equal(expectedIncident))
			})
		})
	})

	Describe("Acknowledge", func() {
		var (
			resp *Response
			err  error
		)

		Context("with a successful, non-empty response", func() {
			BeforeEach(func() {
				path, _ := regexp.Compile("/incidents/\\w+/acknowledge")
				env.Server.RouteToHandler(PUT, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, nil),
				))

				resp, err = env.Client.Incidents.Acknowledge("id", nil)
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

			It("should return a response with the correct status code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Reassign", func() {
		var (
			resp *Response
			err  error
		)

		Context("with a successful, non-empty response", func() {
			BeforeEach(func() {
				path, _ := regexp.Compile("/incidents/\\w+/reassign")
				env.Server.RouteToHandler(PUT, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, nil),
				))

				resp, err = env.Client.Incidents.Reassign("id", nil)
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

			It("should return a response with the correct status code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Resolve", func() {
		var (
			resp *Response
			err  error
		)

		Context("with a successful, non-empty response", func() {
			BeforeEach(func() {
				path, _ := regexp.Compile("/incidents/\\w+/resolve")
				env.Server.RouteToHandler(PUT, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, nil),
				))

				resp, err = env.Client.Incidents.Resolve("id", nil)
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

			It("should return a response with the correct status code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Snooze", func() {
		var (
			resp *Response
			err  error
		)

		Context("with a successful, non-empty response", func() {
			BeforeEach(func() {
				path, _ := regexp.Compile("/incidents/\\w+/snooze")
				env.Server.RouteToHandler(PUT, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, nil),
				))

				resp, err = env.Client.Incidents.Snooze("id", nil)
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

			It("should return a response with the correct status code", func() {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})
