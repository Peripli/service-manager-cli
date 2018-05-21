package smclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSmClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Service Manager Client test", func() {
	var client Client
	var responseStatusCode int
	var responseBody []byte
	var smServer *httptest.Server

	platform := &types.Platform{
		ID:          "1234",
		Name:        "cfeu10",
		Type:        "cf",
		Description: "Test platform",
	}

	broker := &types.Broker{
		Name:        "test broker",
		URL:         "http://test-url.com",
		Credentials: types.Credentials{Basic: types.Basic{User: "test user", Password: "test password"}},
	}

	createSMHandler := func() http.HandlerFunc {
		return func(response http.ResponseWriter, req *http.Request) {
			response.WriteHeader(responseStatusCode)
			response.Write([]byte(responseBody))
		}
	}

	BeforeEach(func() {
		smServer = httptest.NewServer(createSMHandler())
		clientConfig := &ClientConfig{smServer.URL, "admin", "token"}
		client = NewClient(clientConfig)
	})

	Describe("Register platform", func() {
		Context("When valid platform is being registered", func() {
			It("should register successfully", func() {
				responseStatusCode = http.StatusCreated
				responseBody, _ = json.Marshal(platform)

				responsePlatform, err := client.RegisterPlatform(platform)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responsePlatform).To(Equal(platform))
			})
		})

		Context("When invalid platform is returned by SM", func() {
			It("should return error", func() {
				responseBody, _ = json.Marshal(struct {
					Name bool `json:"name"`
				}{
					Name: true,
				})
				responseStatusCode = http.StatusCreated

				responsePlatform, err := client.RegisterPlatform(platform)

				Expect(err).Should(HaveOccurred())
				Expect(responsePlatform).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is successful", func() {
				It("should return error with status code", func() {
					responseBody, _ = json.Marshal(platform)
					responseStatusCode = http.StatusOK

					responsePlatform, err := client.RegisterPlatform(platform)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{StatusCode: responseStatusCode}))
					Expect(responsePlatform).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				It("should return error with url and description", func() {
					responseBody = []byte(`{ "description": "error"}`)
					responseStatusCode = http.StatusBadRequest

					responsePlatform, err := client.RegisterPlatform(platform)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + "/v1/platforms", Description: "error", StatusCode: responseStatusCode}))
					Expect(responsePlatform).To(BeNil())
				})
			})

			Context("And response body is invalid", func() {
				It("should return error without url and description if invalid response body", func() {
					responseStatusCode = http.StatusBadRequest
					responseBody = []byte(`{ "description": description", "error": "error"}`)

					responsePlatform, err := client.RegisterPlatform(platform)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + "/v1/platforms", StatusCode: responseStatusCode}))
					Expect(responsePlatform).To(BeNil())
				})
			})
		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				clientConfig := &ClientConfig{"invalidURL", "admin", "token"}
				client = NewClient(clientConfig)
				_, err := client.RegisterPlatform(platform)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Register broker", func() {
		Context("When valid broker is being registered", func() {
			It("should register successfully", func() {
				responseStatusCode = http.StatusCreated
				responseBody, _ = json.Marshal(broker)

				responseBroker, err := client.RegisterBroker(broker)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responseBroker).To(Equal(broker))
			})
		})

		Context("When invalid broker is being returned by SM", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated
				responseBody, _ = json.Marshal(struct {
					Name bool `json:"name"`
				}{
					Name: true,
				})

				responseBroker, err := client.RegisterBroker(broker)

				Expect(err).Should(HaveOccurred())
				Expect(responseBroker).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is unsuccessful", func() {
				It("should return error with status code", func() {
					responseStatusCode = http.StatusOK
					responseBody, _ = json.Marshal(broker)

					responseBroker, err := client.RegisterBroker(broker)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{StatusCode: responseStatusCode}))
					Expect(responseBroker).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				It("should return error with url and description", func() {
					responseStatusCode = http.StatusBadRequest
					responseBody = []byte(`{ "description": "description", "error": "error"}`)

					responseBroker, err := client.RegisterBroker(broker)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + "/v1/service_brokers", Description: "description", ErrorMessage: "error", StatusCode: responseStatusCode}))
					Expect(responseBroker).To(BeNil())
				})
			})

			Context("And invalid response body", func() {
				It("should return error without url and description if invalid response body", func() {
					responseStatusCode = http.StatusBadRequest
					responseBody = []byte(`{ "description": description", "error": "error"}`)

					responseBroker, err := client.RegisterBroker(broker)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + "/v1/service_brokers", StatusCode: responseStatusCode}))
					Expect(responseBroker).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				clientConfig := &ClientConfig{"invalidURL", "admin", "token"}
				client = NewClient(clientConfig)
				_, err := client.RegisterBroker(broker)

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
