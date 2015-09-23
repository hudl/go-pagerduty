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
	escalationPolicyListJSON = `{"escalation_policies": [` + escalationPolicyJSON + `]}`
	escalationPolicyGetJSON  = `{"escalation_policy": ` + escalationPolicyJSON + `}`
	escalationPolicyJSON     = `{
		"id": "id",
		"name": "name",
		"escalation_rules": [{
			"id": "id",
			"escalation_delay_in_minutes": 1,
			"targets": [{
				"id": "id",
				"name": "name",
				"type": "type",
				"email": "email",
				"time_zone": "UTC",
				"color": "color"
			}]
		}]
	}`
)

var _ = Describe("EscalationPolicies", func() {
	var (
		env                      *TestEnvironment
		expectedEscalationPolicy EscalationPolicy
	)

	json.Unmarshal([]byte(escalationPolicyJSON), &expectedEscalationPolicy)

	BeforeEach(func() { env = NewTestEnvironment() })
	AfterEach(func() { env.Server.Close() })

	Describe("List", func() {
		Context("with a successful, non-empty response", func() {
			var (
				policies []EscalationPolicy
				resp     *Response
				err      error
			)

			BeforeEach(func() {
				env.Server.RouteToHandler(GET, "/escalation_policies", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, escalationPolicyListJSON),
				))

				policies, resp, err = env.Client.EscalationPolicies.List(nil)
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

			It("should return the expected escalation policies", func() {
				Expect(policies).NotTo(BeNil())
				Expect(policies).NotTo(BeEmpty())
				Expect(policies[0]).To(Equal(expectedEscalationPolicy))
			})
		})
	})

	Describe("Get", func() {
		Context("with a successful, non-empty response", func() {
			var (
				policy *EscalationPolicy
				resp   *Response
				err    error
			)

			BeforeEach(func() {
				path, _ := regexp.Compile("/escalation_policies/\\w+")
				env.Server.RouteToHandler(GET, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, escalationPolicyGetJSON),
				))

				policy, resp, err = env.Client.EscalationPolicies.Get("id")
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

			It("should return the expected escalation policy", func() {
				Expect(policy).NotTo(BeNil())
				Expect(policy).To(Equal(&expectedEscalationPolicy))
			})
		})
	})
})
