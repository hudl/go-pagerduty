package pagerduty_test

import (
	. "github.com/hudl/go-pagerduty/pagerduty"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"net/http"
)

const alertListJSON = `{ "alerts": [{ "id": "id", "type": "type" }] }`

var _ = Describe("Alerts", func() {
	var env *TestEnvironment

	Describe("List", func() {
		BeforeEach(func() { env = NewTestEnvironment() })
		AfterEach(func() { env.Server.Close() })

		Context("with a successful, non-empty response", func() {
			var (
				alerts []Alert
				resp   *Response
				err    error
			)

			BeforeEach(func() {
				env.Server.RouteToHandler(GET, "/alerts", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, alertListJSON),
				))

				alerts, resp, err = env.Client.Alerts.List(nil)
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return a non-nil response", func() {
				Expect(resp).NotTo(BeNil())
			})

			It("should return the expected alerts", func() {
				Expect(alerts).NotTo(BeNil())
				Expect(alerts).NotTo(BeEmpty())
				Expect(alerts[0]).To(Equal(Alert{
					ID:   String("id"),
					Type: String("type"),
				}))
			})
		})
	})
})
