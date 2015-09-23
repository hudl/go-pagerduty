package pagerduty_test

import (
	. "github.com/hudl/go-pagerduty/pagerduty"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"encoding/json"
	"net/http"
	"regexp"
)

const (
	serviceListJSON = `{ "services": [` + serviceJSON + `]}`
	serviceJSON     = `{
		"id": "id",
		"name": "name",
		"service_url": "service_url",
		"service_key": "service_key",
		"auto_resolve_timeout": 1,
		"acknowledgement_timeout": 1,
		"created_at": null,
		"status": "status",
		"last_incident_timestamp": null,
		"email_incident_creation": "email_incident_creation",
		"incident_counts": {
			"triggered": 1,
			"acknowledged": 1,
			"resolved": 1,
			"total": 3
		},
		"email_filter_mode": "email_filter_mode",
		"type": "type",
		"description": "description"
	}`
)

var _ = Describe("Services", func() {
	var (
		env             *TestEnvironment
		expectedService Service
	)

	json.Unmarshal([]byte(serviceJSON), &expectedService)

	BeforeEach(func() { env = NewTestEnvironment() })
	AfterEach(func() { env.Server.Close() })

	Describe("List", func() {
		Context("with a successful, non-empty response", func() {
			var (
				services []Service
				resp     *Response
				err      error
			)

			BeforeEach(func() {
				env.Server.RouteToHandler(GET, "/services", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, serviceListJSON),
				))

				services, resp, err = env.Client.Services.List(nil)
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

			It("should return the expected service", func() {
				Expect(services).NotTo(BeNil())
				Expect(services).NotTo(BeEmpty())
				Expect(services[0]).To(Equal(expectedService))
			})
		})
	})

	Describe("Get", func() {
		Context("with a successful, non-empty response", func() {
			var (
				service *Service
				resp    *Response
				err     error
			)

			BeforeEach(func() {
				path, _ := regexp.Compile("/services/\\w+")
				env.Server.RouteToHandler(GET, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, serviceJSON),
				))

				service, resp, err = env.Client.Services.Get("id", nil)
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

			It("should return the expected service", func() {
				Expect(service).NotTo(BeNil())
				Expect(service).To(Equal(&expectedService))
			})
		})
	})
})
