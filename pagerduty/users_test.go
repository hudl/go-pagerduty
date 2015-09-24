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
	userListJSON = `{ "users": [` + userJSON + `]}`
	userGetJSON  = `{ "user": ` + userJSON + `}`
	userJSON     = `{
		"id": "id",
		"name": "name",
		"email": "email",
		"time_zone": "Eastern Time (US & Canada)",
		"color": "color",
		"role": "role",
		"avatar_url": "avatar_url",
		"user_url": "user_url",
		"invitation_sent": false
	}`
)

var _ = Describe("Users", func() {
	var (
		env          *TestEnvironment
		expectedUser User
	)

	json.Unmarshal([]byte(userJSON), &expectedUser)

	BeforeEach(func() { env = NewTestEnvironment() })
	AfterEach(func() { env.Server.Close() })

	Describe("List", func() {
		Context("with a successful, non-empty response", func() {
			var (
				users []User
				resp  *Response
				err   error
			)

			BeforeEach(func() {
				env.Server.RouteToHandler(GET, "/users", ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, userListJSON),
				))

				users, resp, err = env.Client.Users.List(nil)
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

			It("should return the expected user", func() {
				Expect(users).NotTo(BeNil())
				Expect(users).NotTo(BeEmpty())
				Expect(users[0]).To(Equal(expectedUser))
			})
		})
	})

	Describe("Get", func() {
		Context("with a successful, non-empty response", func() {
			var (
				user *User
				resp *Response
				err  error
			)

			BeforeEach(func() {
				path, _ := regexp.Compile("/users/\\w+")
				env.Server.RouteToHandler(GET, path, ghttp.CombineHandlers(
					verifyHeaderHandler,
					ghttp.RespondWith(http.StatusOK, userGetJSON),
				))

				user, resp, err = env.Client.Users.Get("id", nil)
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

			It("should return the expected user", func() {
				Expect(user).NotTo(BeNil())
				Expect(user).To(Equal(&expectedUser))
			})
		})
	})
})
