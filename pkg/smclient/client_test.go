package smclient

import (
	"encoding/json"
	"github.com/Peripli/service-manager/pkg/web"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeAuthClient struct {
	AccessToken string
}

func (c *FakeAuthClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	return http.DefaultClient.Do(req)
}

type HandlerDetails struct {
	Method string
	Path string
	ResponseBody []byte
	ResponseStatusCode int
}

func TestSmClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Service Manager Client test", func() {
	var client Client
	var handlerDetails []HandlerDetails
	var validToken = "valid-token"
	var smServer *httptest.Server

	platform := &types.Platform{
		ID:          "1234",
		Name:        "cfeu10",
		Type:        "cf",
		Description: "Test platform",
	}

	broker := &types.Broker{
		Name:        "test-broker",
		URL:         "http://test-url.com",
		Credentials: &types.Credentials{Basic: types.Basic{User: "test user", Password: "test password"}},
	}

	initialOffering := &types.ServiceOffering{
		ID: "offeringID",
		Name: "initial-offering",
		Description: "Some description",
		BrokerID: "id",
	}

	plan := &types.ServicePlan{
		ID: "planID",
		Name: "plan-1",
		Description: "Sample Plan",
		ServiceOfferingID: "offeringID",
	}

	resultOffering := &types.ServiceOffering{
		ID: "offeringID",
		Name: "initial-offering",
		Description: "Some description",
		Plans: []types.ServicePlan{*plan},
		BrokerID: "id",
		BrokerName: "test-broker",
	}



	createSMHandler := func() http.Handler {
		mux := http.NewServeMux()
		for i := range handlerDetails {
			v := handlerDetails[i]
			mux.HandleFunc(v.Path, func(response http.ResponseWriter, req *http.Request) {
				if v.Method != req.Method {
					return
				}
				authorization := req.Header.Get("Authorization")
				if authorization != "Bearer "+validToken {
					response.WriteHeader(http.StatusUnauthorized)
					response.Write([]byte(""))
					return
				}
				response.WriteHeader(v.ResponseStatusCode)
				response.Write(v.ResponseBody)
			})
		}
		return mux
	}

	JustBeforeEach(func() {
		smServer = httptest.NewServer(createSMHandler())
		fakeAuthClient := &FakeAuthClient{AccessToken: validToken}
		client = NewClient(fakeAuthClient, smServer.URL)
	})

	Describe("Test failing client authentication", func() {
		Context("When wrong token is used", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.BrokersURL},
				}
			})
			It("should fail to authentication", func() {
				client = NewClient(http.DefaultClient, smServer.URL)
				_, err := client.ListBrokers()

				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.BrokersURL, StatusCode: http.StatusUnauthorized}))
			})
		})
	})

	Describe("Register platform", func() {
		Context("When valid platform is being registered", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(platform)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.PlatformsURL , ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responsePlatform, err := client.RegisterPlatform(platform)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responsePlatform).To(Equal(platform))
			})
		})

		Context("When invalid platform is returned by SM", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(struct {
					Name bool `json:"name"`
				}{
					Name: true,
				})
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should return error", func() {
				responsePlatform, err := client.RegisterPlatform(platform)

				Expect(err).Should(HaveOccurred())
				Expect(responsePlatform).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is successful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(platform)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
					}
				})
				It("should return error with status code", func() {
					responsePlatform, err := client.RegisterPlatform(platform)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{StatusCode: handlerDetails[0].ResponseStatusCode}))
					Expect(responsePlatform).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responsePlatform, err := client.RegisterPlatform(platform)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.PlatformsURL, Description: "error", StatusCode: handlerDetails[0].ResponseStatusCode}))
					Expect(responsePlatform).To(BeNil())
				})
			})

			Context("And response body is invalid", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responsePlatform, err := client.RegisterPlatform(platform)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.PlatformsURL, StatusCode: handlerDetails[0].ResponseStatusCode}))
					Expect(responsePlatform).To(BeNil())
				})
			})
		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = NewClient(http.DefaultClient, "invalidURL")
				_, err := client.RegisterPlatform(platform)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Register broker", func() {
		Context("When valid broker is being registered", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.BrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responseBroker, err := client.RegisterBroker(broker)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responseBroker).To(Equal(broker))
			})
		})

		Context("When invalid broker is being returned by SM", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(struct {
					Name bool `json:"name"`
				}{
					Name: true,
				})
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.BrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should return error", func() {
				responseBroker, err := client.RegisterBroker(broker)

				Expect(err).Should(HaveOccurred())
				Expect(responseBroker).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(broker)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.BrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
					}
				})
				It("should return error with status code", func() {
					responseBroker, err := client.RegisterBroker(broker)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{StatusCode: handlerDetails[0].ResponseStatusCode}))
					Expect(responseBroker).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.BrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responseBroker, err := client.RegisterBroker(broker)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.BrokersURL, Description: "description",
																	ErrorMessage: "error", StatusCode: handlerDetails[0].ResponseStatusCode}))
					Expect(responseBroker).To(BeNil())
				})
			})

			Context("And invalid response body", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.BrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responseBroker, err := client.RegisterBroker(broker)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.BrokersURL, StatusCode: handlerDetails[0].ResponseStatusCode}))
					Expect(responseBroker).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = NewClient(http.DefaultClient, "invalidURL")
				_, err := client.RegisterBroker(broker)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("List brokers", func() {
		Context("when there are brokers registered", func() {
			BeforeEach(func() {
				brokersArray := []types.Broker{*broker}
				brokers := types.Brokers{Brokers: brokersArray}
				responseBody, _ := json.Marshal(brokers)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.BrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all", func() {
				result, err := client.ListBrokers()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Brokers).To(HaveLen(1))
				Expect(result.Brokers[0]).To(Equal(*broker))
			})
		})

		Context("when there are no brokers registered", func() {
			BeforeEach(func() {
				brokersArray := []types.Broker{}
				brokers := types.Brokers{Brokers: brokersArray}
				responseBody, _ := json.Marshal(brokers)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.BrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return empty array", func() {
				result, err := client.ListBrokers()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Brokers).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.BrokersURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListBrokers()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.BrokersURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListBrokers()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusBadRequest, URL: smServer.URL + web.BrokersURL + "?fieldQuery=&labelQuery="}))
			})
		})
	})

	Describe("List platforms", func() {
		Context("when there are platforms registered", func() {
			BeforeEach(func() {
				platformsArray := []types.Platform{*platform}
				platforms := types.Platforms{Platforms: platformsArray}
				responseBody, _ := json.Marshal(platforms)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all", func() {
				result, err := client.ListPlatforms()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Platforms).To(HaveLen(1))
				Expect(result.Platforms[0]).To(Equal(*platform))
			})
		})

		Context("when there are no platforms registered", func() {
			BeforeEach(func() {
				platformsArray := []types.Platform{}
				platforms := types.Platforms{Platforms: platformsArray}
				responseBody, _ := json.Marshal(platforms)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return empty array", func() {
				result, err := client.ListPlatforms()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Platforms).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListPlatforms()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})
		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListPlatforms()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusBadRequest, URL: smServer.URL + web.PlatformsURL + "?fieldQuery="}))
			})
		})
	})

	Describe("List offerings", func() {
		Context("when there are offerings provided", func() {
			BeforeEach(func() {
				offerings := types.ServiceOfferings{ServiceOfferings: []types.ServiceOffering{*initialOffering}}
				offeringResponseBody, _ := json.Marshal(offerings)

				plans := types.ServicePlans{ServicePlans: []types.ServicePlan{*plan}}
				plansResponseBody, _ := json.Marshal(plans)

				brokerResponseBody, _ := json.Marshal(broker)

				handlerDetails = []HandlerDetails {
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseBody: offeringResponseBody, ResponseStatusCode: http.StatusOK},
					{Method: http.MethodGet, Path: web.ServicePlansURL, ResponseBody: plansResponseBody, ResponseStatusCode: http.StatusOK},
					{Method: http.MethodGet, Path: web.BrokersURL + "/", ResponseBody: brokerResponseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all with plans and broker name populated", func() {
				result, err := client.ListOfferings()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.ServiceOfferings).To(HaveLen(1))
				Expect(result.ServiceOfferings[0]).To(Equal(*resultOffering))
			})
		})

		Context("when there are no offerings provided", func() {
			BeforeEach(func() {
				offerings := types.ServiceOfferings{}
				offeringResponseBody, _ := json.Marshal(offerings)

				handlerDetails = []HandlerDetails {
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseBody: offeringResponseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return empty array", func() {
				result, err := client.ListOfferings()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.ServiceOfferings).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails {
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListOfferings()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails {
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListOfferings()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusBadRequest, URL: smServer.URL + web.ServiceOfferingsURL}))
			})
		})

	})

	Describe("Delete brokers", func() {
		Context("when an existing broker is being deleted", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.BrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				err := client.DeleteBroker("id")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.BrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				err := client.DeleteBroker("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Broker not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.BrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				err := client.DeleteBroker("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Broker not found", URL: smServer.URL + web.BrokersURL + "/id", StatusCode: http.StatusNotFound}))
			})
		})
	})

	Describe("Delete platforms", func() {
		Context("when an existing platform is being deleted", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				err := client.DeletePlatform("id")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				err := client.DeletePlatform("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Platform not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				err := client.DeletePlatform("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Platform not found", URL: smServer.URL + web.PlatformsURL + "/id", StatusCode: http.StatusNotFound}))
			})
		})
	})

	Describe("Update brokers", func() {
		Context("when an existing broker is being updated", func() {
			BeforeEach(func() {
				responseBody := []byte(`{
					"id": "1234",
					"name": "broker",
					"broker_url": "http://broker.com"
				}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.BrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				updatedBroker, err := client.UpdateBroker("1234", &types.Broker{Name: "broker"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updatedBroker.Name).To(Equal("broker"))
			})
		})

		Context("when a non-existing broker is being updated", func() {
			BeforeEach(func() {
				responseBody := []byte(`{}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.BrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdateBroker("1234", &types.Broker{Name: "broker"})
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Get info", func() {
		Context("when token issuer is set", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"token_issuer_url": "http://uaa.com"}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.InfoURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should get the right issuer and default token basic auth", func() {
				info, _ := client.GetInfo()
				Expect(info.TokenIssuerURL).To(Equal("http://uaa.com"))
				Expect(info.TokenBasicAuth).To(BeTrue()) // default value
			})
		})

		Context("when token basic auth is set", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"token_issuer_url": "http://uaa.com", "token_basic_auth": false}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.InfoURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should get the right value", func() {
				info, _ := client.GetInfo()
				Expect(info.TokenIssuerURL).To(Equal("http://uaa.com"))
				Expect(info.TokenBasicAuth).To(BeFalse())
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				responseBody := []byte(``)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.InfoURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should get an error", func() {
				_, err := client.GetInfo()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.InfoURL, StatusCode: http.StatusNotFound}))
			})
		})

		Context("when invalid json is returned", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"token_issuer":}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.InfoURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should get an error", func() {
				_, err := client.GetInfo()
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
