package pagerduty_test

import (
	. "github.com/hudl/go-pagerduty/pagerduty"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"encoding/json"
	"net/http"
)

const (
	eventsAPIURL = "/generic/2010-04-15/create_event.json"

	eventSuccessResponseJSON = `{
		"status": "status",
		"message": "message",
		"incident_key": "incident_key"
	}`

	eventErrorResponseJSON = `{
		"status": "status",
		"message": "message",
		"errors": [
			"error message"
		]
	}`
)

func verifyEventType(eventType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&event)
		Expect(err).NotTo(HaveOccurred())

		Expect(event).To(HaveKey("event_type"))
		Expect(event["event_type"]).To(BeAssignableToTypeOf(*new(string)))
		Expect(event["event_type"].(string)).To(Equal(eventType))
	}
}

var _ = Describe("Events", func() {
	var (
		env  *TestEnvironment
		resp *EventResponse
		err  error
	)

	BeforeEach(func() { env = NewTestEnvironment() })
	AfterEach(func() { env.Server.Close() })

	Describe("Acknowledge", func() {
		Context("with a nil event", func() {
			BeforeEach(func() {
				resp, err = env.Client.Events.Acknowledge(nil)
			})

			It("should return an error", func() {
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with a non-nil event", func() {
			BeforeEach(func() {
				env.Server.RouteToHandler(POST, eventsAPIURL, ghttp.CombineHandlers(
					verifyContentHeaderHandler,
					verifyEventType(EventTypeAcknowledge),
					ghttp.RespondWith(http.StatusOK, eventSuccessResponseJSON),
				))

				resp, err = env.Client.Events.Acknowledge(new(Event))
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the expected event response", func() {
				Expect(resp).NotTo(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(new(EventResponse)))
				Expect(resp).To(Equal(&EventResponse{
					Response:    resp.Response,
					Status:      "status",
					Message:     "message",
					IncidentKey: "incident_key",
				}))
			})
		})

		Context("when an error occurs on the server", func() {
			BeforeEach(func() {
				env.Server.RouteToHandler(POST, eventsAPIURL, ghttp.CombineHandlers(
					verifyContentHeaderHandler,
					verifyEventType(EventTypeAcknowledge),
					ghttp.RespondWith(http.StatusBadRequest, eventErrorResponseJSON),
				))

				resp, err = env.Client.Events.Acknowledge(new(Event))
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the expected event response", func() {
				Expect(resp).NotTo(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(new(EventResponse)))
				Expect(resp).To(Equal(&EventResponse{
					Response: resp.Response,
					Status:   "status",
					Message:  "message",
					Errors:   []string{"error message"},
				}))
			})
		})
	})

	Describe("Resolve", func() {
		Context("with a nil event", func() {
			BeforeEach(func() {
				resp, err = env.Client.Events.Resolve(nil)
			})

			It("should return an error", func() {
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with a non-nil event", func() {
			BeforeEach(func() {
				env.Server.RouteToHandler(POST, eventsAPIURL, ghttp.CombineHandlers(
					verifyContentHeaderHandler,
					verifyEventType(EventTypeResolve),
					ghttp.RespondWith(http.StatusOK, eventSuccessResponseJSON),
				))

				resp, err = env.Client.Events.Resolve(new(Event))
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the expected event response", func() {
				Expect(resp).NotTo(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(new(EventResponse)))
				Expect(resp).To(Equal(&EventResponse{
					Response:    resp.Response,
					Status:      "status",
					Message:     "message",
					IncidentKey: "incident_key",
				}))
			})
		})

		Context("when an error occurs on the server", func() {
			BeforeEach(func() {
				env.Server.RouteToHandler(POST, eventsAPIURL, ghttp.CombineHandlers(
					verifyContentHeaderHandler,
					verifyEventType(EventTypeResolve),
					ghttp.RespondWith(http.StatusBadRequest, eventErrorResponseJSON),
				))

				resp, err = env.Client.Events.Resolve(new(Event))
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the expected event response", func() {
				Expect(resp).NotTo(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(new(EventResponse)))
				Expect(resp).To(Equal(&EventResponse{
					Response: resp.Response,
					Status:   "status",
					Message:  "message",
					Errors:   []string{"error message"},
				}))
			})
		})
	})

	Describe("Trigger", func() {
		Context("with a nil event", func() {
			BeforeEach(func() {
				resp, err = env.Client.Events.Trigger(nil)
			})

			It("should return an error", func() {
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with a non-nil event", func() {
			BeforeEach(func() {
				env.Server.RouteToHandler(POST, eventsAPIURL, ghttp.CombineHandlers(
					verifyContentHeaderHandler,
					verifyEventType(EventTypeTrigger),
					ghttp.RespondWith(http.StatusOK, eventSuccessResponseJSON),
				))

				resp, err = env.Client.Events.Trigger(new(Event))
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the expected event response", func() {
				Expect(resp).NotTo(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(new(EventResponse)))
				Expect(resp).To(Equal(&EventResponse{
					Response:    resp.Response,
					Status:      "status",
					Message:     "message",
					IncidentKey: "incident_key",
				}))
			})
		})

		Context("when an error occurs on the server", func() {
			BeforeEach(func() {
				env.Server.RouteToHandler(POST, eventsAPIURL, ghttp.CombineHandlers(
					verifyContentHeaderHandler,
					verifyEventType(EventTypeTrigger),
					ghttp.RespondWith(http.StatusBadRequest, eventErrorResponseJSON),
				))

				resp, err = env.Client.Events.Trigger(new(Event))
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the expected event response", func() {
				Expect(resp).NotTo(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(new(EventResponse)))
				Expect(resp).To(Equal(&EventResponse{
					Response: resp.Response,
					Status:   "status",
					Message:  "message",
					Errors:   []string{"error message"},
				}))
			})
		})
	})
})
