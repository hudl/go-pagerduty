package pagerduty_test

import (
	. "github.com/hudl/go-pagerduty/pagerduty"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
)

const (
	scheduleListJSON = `{ "schedules": [` + scheduleJSON + `]}`
	scheduleGetJSON  = `{ "schedule": ` + scheduleJSON + `}`
	scheduleJSON     = `{
		"id": "id",
		"name": "name",
		"time_zone": "Eastern Time (US & Canada)",
		"today": "2006-01-02",
		"escalation_policies": []
	}`
)

var _ = Describe("Schedules", func() {
	var (
		env              *TestEnvironment
		expectedSchedule Schedule
		expectedUser     User

		resp *Response
		err  error
	)

	json.Unmarshal([]byte(scheduleJSON), &expectedSchedule)
	json.Unmarshal([]byte(userJSON), &expectedUser)

	BeforeEach(func() { env = NewTestEnvironment() })
	AfterEach(func() { env.Server.Close() })

	Describe("List", func() {
		Context("with a valid schedule id and options", func() {
			var schedules []Schedule

			BeforeEach(func() {
				env.Server.RouteToHandler(GET, "/schedules", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, scheduleListJSON),
				))

				schedules, resp, err = env.Client.Schedules.List(nil)
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

			It("should return the expected schedule", func() {
				Expect(schedules).NotTo(BeNil())
				Expect(schedules).NotTo(BeEmpty())
				Expect(schedules[0]).To(Equal(expectedSchedule))
			})
		})
	})

	Describe("Get", func() {
		Context("with a valid schedule id and options", func() {
			var schedule *Schedule

			BeforeEach(func() {
				path, _ := regexp.Compile("/schedules/\\w+")
				env.Server.RouteToHandler(GET, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, scheduleGetJSON),
				))

				schedule, resp, err = env.Client.Schedules.Get("id")
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

			It("should return the expected schedule", func() {
				Expect(schedule).NotTo(BeNil())
				Expect(schedule).To(Equal(&expectedSchedule))
			})
		})
	})

	Describe("Users", func() {
		Context("with a valid schedule id and options", func() {
			var users []User

			BeforeEach(func() {
				values := make(url.Values)
				values.Add("since", testTimeString)
				values.Add("until", testTimeString)

				path, _ := regexp.Compile("/schedules/\\w+/users")
				env.Server.RouteToHandler(GET, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					verifyURLQueryHandler(values),
					ghttp.RespondWith(http.StatusOK, userListJSON),
				))

				users, resp, err = env.Client.Schedules.Users("id", &ScheduleUsersOptions{
					Since: testTime,
					Until: testTime,
				})
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

			It("should return the expected users", func() {
				Expect(users).NotTo(BeNil())
				Expect(users).NotTo(BeEmpty())
				Expect(users[0]).To(Equal(expectedUser))
			})
		})
	})
})
