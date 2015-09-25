package pagerduty_test

import (
	. "github.com/hudl/go-pagerduty/pagerduty"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var _ = Describe("PagerDuty", func() {
	Describe("Creating a new client", func() {
		var client *Client

		Context("without an http client", func() {
			BeforeEach(func() {
				client = NewClient(nil, subdomain, apiKey)
			})

			It("should use the default base URL", func() {
				Expect(client.BaseURL.String()).To(Equal(defaultBaseURL))
			})

			It("should have the correct api key", func() {
				Expect(client.APIKey).To(Equal(apiKey))
			})

			It("should register all services correctly", func() {
				Expect(client.Alerts).NotTo(BeNil())
				Expect(client.EscalationPolicies).NotTo(BeNil())
				Expect(client.Events).NotTo(BeNil())
				Expect(client.Incidents).NotTo(BeNil())
				Expect(client.Schedules).NotTo(BeNil())
				Expect(client.Services).NotTo(BeNil())
				Expect(client.Teams).NotTo(BeNil())
				Expect(client.Users).NotTo(BeNil())
				Expect(client.Webhooks).NotTo(BeNil())
			})
		})
	})

	Describe("Creating a new request", func() {
		var (
			client *Client
			err    error
			req    *http.Request
		)

		BeforeEach(func() {
			client = NewClient(nil, subdomain, apiKey)
		})

		Context("with a valid http method, path and body", func() {
			// test type
			type T struct{ A int }

			BeforeEach(func() {
				req, err = client.NewRequest(GET, "test", &T{A: 0})
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should create a non-nil request", func() {
				Expect(req).NotTo(BeNil())
			})

			It("should set the correct headers", func() {
				Expect(req.Header.Get("Authorization")).To(Equal("Token token=" + apiKey))
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
			})

			It("should expand the relative path", func() {
				Expect(req.URL.String()).To(Equal(defaultBaseURL + "test"))
			})

			It("should encode the body as valid JSON", func() {
				body, _ := ioutil.ReadAll(req.Body)

				Expect(body).To(MatchJSON(`{ "A": 0 }`))
			})
		})

		Context("with invalid JSON", func() {
			// test type with an unsupported json type
			type T struct{ A map[int]interface{} }

			It("should return a JSON unsupported type error", func() {
				_, err = client.NewRequest(GET, "/", &T{})
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(new(json.UnsupportedTypeError)))
			})
		})

		Context("with an invalid realtive path", func() {
			It("should return a URL parse error", func() {
				_, err := client.NewRequest(GET, ":", nil)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(new(url.Error)))
				Expect(err.(*url.Error).Op).To(Equal("parse"))
			})
		})

		Context("with an empty (nil) body", func() {
			BeforeEach(func() {
				req, err = client.NewRequest(GET, "/", nil)
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have an empty (nil) body", func() {
				Expect(req.Body).To(BeNil())
			})
		})
	})

	Describe("Performing a request", func() {
		var (
			body       interface{}
			env        *TestEnvironment
			err        error
			req        *http.Request
			resp       *Response
			statusCode int
		)

		JustBeforeEach(func() {
			env = NewTestEnvironment()

			env.Server.RouteToHandler(GET, "/", ghttp.CombineHandlers(
				verifyHeaderHandler,
				ghttp.RespondWith(statusCode, body),
			))

			req, _ = env.Client.NewRequest(GET, "/", nil)
			resp, err = env.Client.Do(req, nil)
		})

		AfterEach(func() {
			// close the test server
			env.Server.Close()
		})

		Context("with a valid request", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				body = `{ "offset": 1, "limit": 1, "total": 1 }`
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should parse and return the correct response body", func() {
				Expect(resp).To(Equal(&Response{
					Response: resp.Response,
					Offset:   1,
					Limit:    1,
					Total:    1,
				}))
			})
		})

		Context("when the response is an http error", func() {
			BeforeEach(func() {
				statusCode = http.StatusBadRequest
				body = nil
			})

			It("should have made a request", func() {
				Expect(env.Server.ReceivedRequests()).To(HaveLen(1))
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Checking a response", func() {
		var (
			body string
			err  error
			resp *http.Response
		)

		JustBeforeEach(func() {
			resp = &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader(body)),
			}
			err = CheckResponse(resp)
		})

		Context("with an error status code", func() {
			BeforeEach(func() {
				body = `{
					"message": "Message",
					"code": 1,
					"errors": []
				}`
			})

			It("should return a PagerDuty error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(new(ErrorResponse)))
				Expect(err).To(Equal(&ErrorResponse{
					Response: resp,
					Message:  "Message",
					Code:     1,
					Errors:   []string{},
				}))
			})
		})

		Context("with no body", func() {
			BeforeEach(func() {
				body = ""
			})

			It("should return a PagerDuty error without an error code and message", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(new(ErrorResponse)))
				Expect(err).To(Equal(&ErrorResponse{
					Response: resp,
				}))
			})
		})
	})

	Describe("Stringifying a PagerDuty error response", func() {
		Context("with a non-nil PagerDuty error", func() {
			It("should return a non-empty string", func() {
				err := &ErrorResponse{}
				Expect(err.Error()).NotTo(BeEmpty())
			})
		})
	})
})
