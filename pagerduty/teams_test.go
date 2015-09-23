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
	teamListJSON = `{ "teams": [` + teamJSON + `]}`
	teamGetJSON  = `{ "team": ` + teamJSON + `}`
	teamJSON     = `{
		"id": "id",
		"name": "name",
		"description": "description"
	}`
)

var _ = Describe("Teams", func() {
	var (
		env          *TestEnvironment
		expectedTeam Team
	)

	json.Unmarshal([]byte(teamJSON), &expectedTeam)

	BeforeEach(func() { env = NewTestEnvironment() })
	AfterEach(func() { env.Server.Close() })

	Describe("List", func() {
		Context("with a successful, non-empty response", func() {
			var (
				teams []Team
				resp  *Response
				err   error
			)

			BeforeEach(func() {
				env.Server.RouteToHandler(GET, "/teams", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, teamListJSON),
				))

				teams, resp, err = env.Client.Teams.List(nil)
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

			It("should return the expected team", func() {
				Expect(teams).NotTo(BeNil())
				Expect(teams).NotTo(BeEmpty())
				Expect(teams[0]).To(Equal(expectedTeam))
			})
		})
	})

	Describe("Get", func() {
		Context("with a successful, non-empty response", func() {
			var (
				team *Team
				resp *Response
				err  error
			)

			BeforeEach(func() {
				path, _ := regexp.Compile("/teams/\\w+")
				env.Server.RouteToHandler(GET, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, teamGetJSON),
				))

				team, resp, err = env.Client.Teams.Get("id")
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

			It("should return the expected team", func() {
				Expect(team).NotTo(BeNil())
				Expect(team).To(Equal(&expectedTeam))
			})
		})
	})
})
