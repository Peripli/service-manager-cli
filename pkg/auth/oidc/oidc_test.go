package oidc

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Peripli/service-manager-cli/pkg/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

func TestAuthStrategy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Service Manager Auth strategy test", func() {
	var authStrategy auth.AuthenticationStrategy
	var authOptions *auth.Options
	var configurationResponseCode int
	var configurationResponseBody []byte
	var responseStatusCode int
	var responseBody []byte
	var uaaServer *httptest.Server

	createUAAHandler := func() http.HandlerFunc {
		return func(response http.ResponseWriter, req *http.Request) {
			response.Header().Add("Content-Type", "application/json")
			if strings.Contains(req.URL.String(), ".well-known/openid-configuration") {
				response.WriteHeader(configurationResponseCode)
				response.Write(configurationResponseBody)
			} else {
				response.WriteHeader(responseStatusCode)
				response.Write([]byte(responseBody))
			}
		}
	}

	BeforeSuite(func() {
		uaaServer = httptest.NewServer(createUAAHandler())
		configurationResponseCode = http.StatusOK
		configurationResponseBody = []byte(`{"token_endpoint": "` + uaaServer.URL + `"}`)

		authStrategy, authOptions, _ = NewOpenIDStrategy(&auth.Options{
			IssuerURL: uaaServer.URL,
		})

		Expect(authOptions).To(Equal(&auth.Options{
			IssuerURL:             uaaServer.URL,
			TokenEndpoint:         uaaServer.URL,
			AuthorizationEndpoint: "",
		}))
	})

	AfterSuite(func() {
		if uaaServer != nil {
			uaaServer.Close()
		}
	})

	BeforeEach(func() {
		configurationResponseCode = http.StatusOK
		configurationResponseBody = []byte(`{"token_endpoint": "` + uaaServer.URL + `"}`)
	})

	Describe("", func() {
		Context("when configuration response is invalid", func() {
			It("should handle wrong response code", func() {
				configurationResponseCode = http.StatusNotFound
				_, _, err := NewOpenIDStrategy(&auth.Options{
					IssuerURL: uaaServer.URL,
				})

				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError("Error occurred while fetching openid configuration: Unexpected status code"))
			})

			It("should handle wrong JSON body", func() {
				configurationResponseCode = http.StatusOK
				configurationResponseBody = []byte(`{"}`)
				_, _, err := NewOpenIDStrategy(&auth.Options{
					IssuerURL: uaaServer.URL,
				})

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("token generation", func() {
		Context("when valid username and password are used", func() {
			It("should issue token", func() {
				responseStatusCode = http.StatusOK
				responseBody = []byte(`{
					"access_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiIzNWFjNDZkZGI0NjQ0YzEyODA1MGI1MDhmOTg3N2M5MSIsInN1YiI6ImYwYmYzNzA1LWMxNWMtNDYxOS1iMzkyLTg2YWYzODRlODkxNiIsInNjb3BlIjpbIm5ldHdvcmsud3JpdGUiLCJjbG91ZF9jb250cm9sbGVyLmFkbWluIiwicm91dGluZy5yb3V0ZXJfZ3JvdXBzLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLndyaXRlIiwibmV0d29yay5hZG1pbiIsImRvcHBsZXIuZmlyZWhvc2UiLCJvcGVuaWQiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMud3JpdGUiLCJzY2ltLnJlYWQiLCJ1YWEudXNlciIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwic2NpbS53cml0ZSJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiZjBiZjM3MDUtYzE1Yy00NjE5LWIzOTItODZhZjM4NGU4OTE2Iiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoiYWRtaW4iLCJlbWFpbCI6ImFkbWluIiwiYXV0aF90aW1lIjoxNTI3NzU3MjMzLCJyZXZfc2lnIjoiYTRiYWI4MTQiLCJpYXQiOjE1Mjc3NTcyMzMsImV4cCI6MTUyNzc1NzgzMywiaXNzIjoiaHR0cHM6Ly91YWEubG9jYWwucGNmZGV2LmlvL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImNsb3VkX2NvbnRyb2xsZXIiLCJzY2ltIiwicGFzc3dvcmQiLCJjZiIsInVhYSIsIm9wZW5pZCIsImRvcHBsZXIiLCJuZXR3b3JrIiwicm91dGluZy5yb3V0ZXJfZ3JvdXBzIl19.Srd_204A3KyHAQ2QibxwxhRm6mwVRRdkJLluiOua6KHmj_x8LLLu6XA9G1e5LNzW_hNqmwxi1fUeFU7NsfUudo46r6pcdfMT0yl7x0qUdizKKZNSkRsoB3BBn1aTBMAgAtc_VBRC8KWCL6Sdy2V0zJ4C-D2nqnYu9vmsK1_tSao",
					"token_type": "bearer",
					"refresh_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiJlNTI2ZDZmNmI4ODk0YjJkOTNhYjI5YTlhY2NmOGNhOS1yIiwic3ViIjoiZjBiZjM3MDUtYzE1Yy00NjE5LWIzOTItODZhZjM4NGU4OTE2Iiwic2NvcGUiOlsibmV0d29yay53cml0ZSIsImNsb3VkX2NvbnRyb2xsZXIuYWRtaW4iLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIud3JpdGUiLCJuZXR3b3JrLmFkbWluIiwiZG9wcGxlci5maXJlaG9zZSIsIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsInNjaW0ucmVhZCIsInVhYS51c2VyIiwiY2xvdWRfY29udHJvbGxlci5yZWFkIiwicGFzc3dvcmQud3JpdGUiLCJzY2ltLndyaXRlIl0sImlhdCI6MTUyNzc1NzIzMywiZXhwIjoxNTMwMzQ5MjMzLCJjaWQiOiJjZiIsImNsaWVudF9pZCI6ImNmIiwiaXNzIjoiaHR0cHM6Ly91YWEubG9jYWwucGNmZGV2LmlvL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9uYW1lIjoiYWRtaW4iLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX2lkIjoiZjBiZjM3MDUtYzE1Yy00NjE5LWIzOTItODZhZjM4NGU4OTE2IiwicmV2X3NpZyI6ImE0YmFiODE0IiwiYXVkIjpbImNsb3VkX2NvbnRyb2xsZXIiLCJzY2ltIiwicGFzc3dvcmQiLCJjZiIsInVhYSIsIm9wZW5pZCIsImRvcHBsZXIiLCJuZXR3b3JrIiwicm91dGluZy5yb3V0ZXJfZ3JvdXBzIl19.fNWVIyrjM7zIf89R1iMwKLNkBwE3Go51OKnnGnpONSsh0KciogcdEN9pYVSZMeb37bDmlc6L-wYUUCSY-ZP4VNm9pZtC-uhIfFy8kT6ZHADpp0IuNbD5AK48NC6yRs8Qgux8OV2UHryxlcMVfCC-EfUUaI6Mcz4JWh1EU7ojesM",
					"expires_in": 599,
					"scope": "network.write cloud_controller.admin routing.router_groups.read cloud_controller.write network.admin doppler.firehose openid routing.router_groups.write scim.read uaa.user cloud_controller.read password.write scim.write",
					"jti": "35ac46ddb4644c128050b508f9877c91"
				}`)

				token, err := authStrategy.Authenticate("admin", "admin")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(token.AccessToken).To(Equal("eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiIzNWFjNDZkZGI0NjQ0YzEyODA1MGI1MDhmOTg3N2M5MSIsInN1YiI6ImYwYmYzNzA1LWMxNWMtNDYxOS1iMzkyLTg2YWYzODRlODkxNiIsInNjb3BlIjpbIm5ldHdvcmsud3JpdGUiLCJjbG91ZF9jb250cm9sbGVyLmFkbWluIiwicm91dGluZy5yb3V0ZXJfZ3JvdXBzLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLndyaXRlIiwibmV0d29yay5hZG1pbiIsImRvcHBsZXIuZmlyZWhvc2UiLCJvcGVuaWQiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMud3JpdGUiLCJzY2ltLnJlYWQiLCJ1YWEudXNlciIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwic2NpbS53cml0ZSJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiZjBiZjM3MDUtYzE1Yy00NjE5LWIzOTItODZhZjM4NGU4OTE2Iiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoiYWRtaW4iLCJlbWFpbCI6ImFkbWluIiwiYXV0aF90aW1lIjoxNTI3NzU3MjMzLCJyZXZfc2lnIjoiYTRiYWI4MTQiLCJpYXQiOjE1Mjc3NTcyMzMsImV4cCI6MTUyNzc1NzgzMywiaXNzIjoiaHR0cHM6Ly91YWEubG9jYWwucGNmZGV2LmlvL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImNsb3VkX2NvbnRyb2xsZXIiLCJzY2ltIiwicGFzc3dvcmQiLCJjZiIsInVhYSIsIm9wZW5pZCIsImRvcHBsZXIiLCJuZXR3b3JrIiwicm91dGluZy5yb3V0ZXJfZ3JvdXBzIl19.Srd_204A3KyHAQ2QibxwxhRm6mwVRRdkJLluiOua6KHmj_x8LLLu6XA9G1e5LNzW_hNqmwxi1fUeFU7NsfUudo46r6pcdfMT0yl7x0qUdizKKZNSkRsoB3BBn1aTBMAgAtc_VBRC8KWCL6Sdy2V0zJ4C-D2nqnYu9vmsK1_tSao"))
			})
		})

		Context("when token response is invalid", func() {
			It("should handle wrong response code", func() {
				errorMsg := `{"error":"missing client_id or client_secret"}`
				responseStatusCode = http.StatusBadRequest
				responseBody = []byte(errorMsg)
				_, err := authStrategy.Authenticate("admin", "admin")

				Expect(err).Should(HaveOccurred())
				Expect(err.(*oauth2.RetrieveError).Error()).To(ContainSubstring(errorMsg))
			})

			It("should handle wrong JSON body", func() {
				responseStatusCode = http.StatusOK
				responseBody = []byte(`{"json":}`)
				_, err := authStrategy.Authenticate("admin", "admin")

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
