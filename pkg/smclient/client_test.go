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
	var validToken = "valid-token"
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
		Credentials: &types.Credentials{Basic: types.Basic{User: "test user", Password: "test password"}},
	}

	createSMHandler := func() http.HandlerFunc {
		return func(response http.ResponseWriter, req *http.Request) {
			authorization := req.Header.Get("Authorization")
			if authorization != "Bearer: "+validToken {
				response.WriteHeader(http.StatusUnauthorized)
				response.Write([]byte(""))
				return
			}
			response.WriteHeader(responseStatusCode)
			response.Write([]byte(responseBody))
		}
	}

	BeforeEach(func() {
		smServer = httptest.NewServer(createSMHandler())
		clientConfig := &ClientConfig{smServer.URL, "admin", "valid-token"}
		client = NewClient(clientConfig)
	})

	Describe("Test failing client authentication", func() {
		Context("When wrong token is used", func() {
			It("should fail to authentication", func() {
				clientConfig := &ClientConfig{smServer.URL, "admin", "invalid-token"}
				client = NewClient(clientConfig)
				_, err := client.ListBrokers()

				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + "/v1/service_brokers", StatusCode: http.StatusUnauthorized}))
			})
		})
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

	Describe("List brokers", func() {
		Context("when there are brokers registered", func() {
			It("should return all", func() {
				responseStatusCode = http.StatusOK

				brokersArray := []types.Broker{*broker}
				brokers := types.Brokers{Brokers: brokersArray}
				responseBody, _ = json.Marshal(brokers)

				result, err := client.ListBrokers()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Brokers).To(HaveLen(1))
				Expect(result.Brokers[0]).To(Equal(*broker))
			})
		})

		Context("when there are no brokers registered", func() {
			It("should return empty array", func() {
				responseStatusCode = http.StatusOK

				brokersArray := []types.Broker{}
				brokers := types.Brokers{Brokers: brokersArray}
				responseBody, _ = json.Marshal(brokers)

				result, err := client.ListBrokers()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Brokers).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			It("should handle status code != 200", func() {
				responseStatusCode = http.StatusCreated

				_, err := client.ListBrokers()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})

			It("should handle status code > 299", func() {
				responseStatusCode = http.StatusBadRequest

				_, err := client.ListBrokers()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusBadRequest, URL: smServer.URL + "/v1/service_brokers"}))
			})
		})
	})

	Describe("List platforms", func() {
		Context("when there are platforms registered", func() {
			It("should return all", func() {
				responseStatusCode = http.StatusOK

				platformsArray := []types.Platform{*platform}
				platforms := types.Platforms{Platforms: platformsArray}
				responseBody, _ = json.Marshal(platforms)

				result, err := client.ListPlatforms()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Platforms).To(HaveLen(1))
				Expect(result.Platforms[0]).To(Equal(*platform))
			})
		})

		Context("when there are no platforms registered", func() {
			It("should return empty array", func() {
				responseStatusCode = http.StatusOK

				platformsArray := []types.Platform{}
				platforms := types.Platforms{Platforms: platformsArray}
				responseBody, _ = json.Marshal(platforms)

				result, err := client.ListPlatforms()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Platforms).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			It("should handle status code != 200", func() {
				responseStatusCode = http.StatusCreated

				_, err := client.ListPlatforms()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})

			It("should handle status code > 299", func() {
				responseStatusCode = http.StatusBadRequest

				_, err := client.ListPlatforms()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusBadRequest, URL: smServer.URL + "/v1/platforms"}))
			})
		})
	})

	Describe("Delete brokers", func() {
		Context("when an existing broker is being deleted", func() {
			It("should be successfully removed", func() {
				responseStatusCode = http.StatusOK
				responseBody = []byte("{}")

				err := client.DeleteBroker("id")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			It("should handle error", func() {
				responseStatusCode = http.StatusCreated
				responseBody = []byte("{}")

				err := client.DeleteBroker("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})

		Context("when service manager returns a status code not found", func() {
			It("should handle error", func() {
				responseStatusCode = http.StatusNotFound
				responseBody = []byte(`{ "description": "Broker not found" }`)

				err := client.DeleteBroker("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Broker not found", URL: smServer.URL + "/v1/service_brokers/id", StatusCode: http.StatusNotFound}))
			})
		})
	})

	Describe("Delete platforms", func() {
		Context("when an existing platform is being deleted", func() {
			It("should be successfully removed", func() {
				responseStatusCode = http.StatusOK
				responseBody = []byte("{}")

				err := client.DeletePlatform("id")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			It("should handle error", func() {
				responseStatusCode = http.StatusCreated
				responseBody = []byte("{}")

				err := client.DeletePlatform("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})

		Context("when service manager returns a status code not found", func() {
			It("should handle error", func() {
				responseStatusCode = http.StatusNotFound
				responseBody = []byte(`{ "description": "Platform not found" }`)

				err := client.DeletePlatform("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Platform not found", URL: smServer.URL + "/v1/platforms/id", StatusCode: http.StatusNotFound}))
			})
		})
	})

	Describe("Update brokers", func() {
		Context("when an existing broker is being updated", func() {
			It("should be successfully removed", func() {
				responseStatusCode = http.StatusOK
				responseBody = []byte(`{
					"id": "1234",
					"name": "broker",
					"broker_url": "http://broker.com"
				}`)

				updatedBroker, err := client.UpdateBroker("1234", &types.Broker{Name: "broker"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updatedBroker.Name).To(Equal("broker"))
			})
		})

		Context("when a non-existing broker is being updated", func() {
			It("should handle error", func() {
				responseStatusCode = http.StatusNotFound
				responseBody = []byte(`{}`)

				_, err := client.UpdateBroker("1234", &types.Broker{Name: "broker"})
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Get info", func() {
		Context("when token issuer is set", func() {
			It("should get the right issuer", func() {
				responseStatusCode = http.StatusOK
				responseBody = []byte(`{"token_issuer_url": "http://uaa.com"}`)

				info, _ := client.GetInfo()
				Expect(info.TokenIssuerURL).To(Equal("http://uaa.com"))
			})
		})

		Context("when invalid status code is returned", func() {
			It("should get an error", func() {
				responseStatusCode = http.StatusNotFound
				responseBody = []byte(``)

				_, err := client.GetInfo()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + "/v1/info", StatusCode: http.StatusNotFound}))
			})
		})

		Context("when invalid json is returned", func() {
			It("should get an error", func() {
				responseStatusCode = http.StatusOK
				responseBody = []byte(`{"token_issuer":}`)

				_, err := client.GetInfo()
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
