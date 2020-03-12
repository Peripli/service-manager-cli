package smclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	cliquery "github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager/pkg/web"

	"github.com/Peripli/service-manager-cli/pkg/types"
	smtypes "github.com/Peripli/service-manager/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeAuthClient struct {
	AccessToken string
	requestURI  string
}

func (c *FakeAuthClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	c.requestURI = req.URL.RequestURI()
	return http.DefaultClient.Do(req)
}

type HandlerDetails struct {
	Method             string
	Path               string
	ResponseBody       []byte
	ResponseStatusCode int
	Headers            map[string]string
}

func TestSMClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Service Manager Client test", func() {
	var client Client
	var handlerDetails []HandlerDetails
	var validToken = "valid-token"
	var invalidToken = "invalid-token"
	var smServer *httptest.Server
	var fakeAuthClient *FakeAuthClient
	var params *cliquery.Parameters

	platform := &types.Platform{
		ID:          "platformID",
		Name:        "cfeu10",
		Type:        "cf",
		Description: "Test platform",
	}

	broker := &types.Broker{
		ID:          "broker-id",
		Name:        "test-broker",
		URL:         "http://test-url.com",
		Credentials: &types.Credentials{Basic: types.Basic{User: "test user", Password: "test password"}},
	}

	operation := &types.Operation{
		ID:         "operation-id",
		Type:       "create",
		State:      "failed",
		ResourceID: "broker-id",
	}

	initialOffering := &types.ServiceOffering{
		ID:          "offeringID",
		Name:        "initial-offering",
		Description: "Some description",
		BrokerID:    "id",
	}

	plan := &types.ServicePlan{
		ID:                "planID",
		Name:              "plan-1",
		Description:       "Sample Plan",
		ServiceOfferingID: "offeringID",
	}

	resultOffering := &types.ServiceOffering{
		ID:          "offeringID",
		Name:        "initial-offering",
		Description: "Some description",
		Plans:       []types.ServicePlan{*plan},
		BrokerID:    "id",
		BrokerName:  "test-broker",
	}

	visibility := &types.Visibility{
		ID:            "visibilityID",
		PlatformID:    "platformID",
		ServicePlanID: "planID",
	}

	instance := &types.ServiceInstance{
		ID:            "instanceID",
		Name:          "instance1",
		ServicePlanID: "service_plan_id",
		PlatformID:    "platform_id",
	}

	binding := &types.ServiceBinding{
		ID:                "instanceID",
		Name:              "instance1",
		ServiceInstanceID: "service_instance_id",
	}

	labelChanges := &types.LabelChanges{
		LabelChanges: []*smtypes.LabelChange{
			{Key: "key", Operation: smtypes.LabelOperation("add"), Values: []string{"val1", "val2"}},
		},
	}

	createSMHandler := func() http.Handler {
		mux := http.NewServeMux()
		for i := range handlerDetails {
			v := handlerDetails[i]
			mux.HandleFunc(v.Path, func(response http.ResponseWriter, req *http.Request) {
				if v.Method != req.Method {
					return
				}
				for key, value := range v.Headers {
					response.Header().Set(key, value)
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

	verifyErrorMsg := func(errorMsg, path string, body []byte, statusCode int) {
		Expect(errorMsg).To(ContainSubstring(buildURL(smServer.URL+path, params)))
		Expect(errorMsg).To(ContainSubstring(string(body)))
		Expect(errorMsg).To(ContainSubstring(fmt.Sprintf("StatusCode: %d", statusCode)))
	}

	BeforeEach(func() {
		params = &cliquery.Parameters{
			GeneralParams: []string{"key=value"},
		}
	})

	AfterEach(func() {
		Expect(fakeAuthClient.requestURI).Should(ContainSubstring("key=value"), fmt.Sprintf("Request URI %s should contain ?key=value", fakeAuthClient.requestURI))
	})

	JustBeforeEach(func() {
		smServer = httptest.NewServer(createSMHandler())
		fakeAuthClient = &FakeAuthClient{AccessToken: validToken}
		client = NewClient(fakeAuthClient, smServer.URL)
	})

	Describe("Get Info", func() {
		BeforeEach(func() {
			responseBody := []byte(`{"token_issuer_url": "http://uaa.com"}`)
			handlerDetails = []HandlerDetails{
				{Method: http.MethodGet, Path: web.InfoURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
			}
		})

		Context("with general parameter", func() {
			It("should make request with these parameters", func() {
				info, err := client.GetInfo(params)
				Expect(err).ToNot(HaveOccurred())
				Expect(info).ToNot(BeNil())
			})
		})
	})

	Describe("Test failing client authentication", func() {
		Context("When wrong token is used", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL},
				}
			})
			It("should fail to authentication", func() {
				fakeAuthClient.AccessToken = invalidToken
				client = NewClient(fakeAuthClient, smServer.URL)
				_, err := client.ListBrokers(params)

				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, []byte{}, http.StatusUnauthorized)
			})
		})
	})

	Describe("Register platform", func() {
		Context("When valid platform is being registered", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(platform)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responsePlatform, err := client.RegisterPlatform(platform, params)

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
				responsePlatform, err := client.RegisterPlatform(platform, params)

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
					responsePlatform, err := client.RegisterPlatform(platform, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
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
					responsePlatform, err := client.RegisterPlatform(platform, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
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
					responsePlatform, err := client.RegisterPlatform(platform, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responsePlatform).To(BeNil())
				})
			})
		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = NewClient(fakeAuthClient, "invalidURL")
				_, err := client.RegisterPlatform(platform, params)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Register Visibility", func() {
		Context("When valid visibility is being registered", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(visibility)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responseVisibility, err := client.RegisterVisibility(visibility, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responseVisibility).To(Equal(visibility))
			})
		})

		Context("When invalid visibility is returned by SM", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(struct {
					ID bool `json:"id"`
				}{
					ID: true,
				})
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})

			It("should return error", func() {
				responseVisibility, err := client.RegisterVisibility(visibility, params)

				Expect(err).Should(HaveOccurred())
				Expect(responseVisibility).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is successful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(visibility)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
					}
				})
				It("should return error with status code", func() {
					responseVisibility, err := client.RegisterVisibility(visibility, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseVisibility).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responseVisibility, err := client.RegisterVisibility(visibility, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseVisibility).To(BeNil())
				})
			})

			Context("And response body is invalid", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responseVisibility, err := client.RegisterVisibility(visibility, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseVisibility).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = NewClient(fakeAuthClient, "invalidURL")
				_, err := client.RegisterVisibility(visibility, params)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Register broker", func() {
		Context("When valid broker is being registered synchronously", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responseBroker, location, err := client.RegisterBroker(broker, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				Expect(responseBroker).To(Equal(broker))
			})
		})

		Context("When valid broker is being registered asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "test-location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should receive operation location", func() {
				responseBroker, location, err := client.RegisterBroker(broker, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
				Expect(responseBroker).To(BeNil())
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
					{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should return error", func() {
				responseBroker, location, err := client.RegisterBroker(broker, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				Expect(responseBroker).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(broker)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
					}
				})
				It("should return error with status code", func() {
					responseBroker, location, err := client.RegisterBroker(broker, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseBroker).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responseBroker, location, err := client.RegisterBroker(broker, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseBroker).To(BeNil())
				})
			})

			Context("And invalid response body", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responseBroker, location, err := client.RegisterBroker(broker, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())

					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseBroker).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = NewClient(fakeAuthClient, "invalidURL")
				_, location, err := client.RegisterBroker(broker, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
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
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all", func() {
				result, err := client.ListBrokers(params)
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
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return an empty array", func() {
				result, err := client.ListBrokers(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Brokers).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListBrokers(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListBrokers(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

			})
		})
	})

	Describe("Get broker", func() {
		Context("when there is broker registered with this id", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return it", func() {
				result, err := client.GetBrokerByID(broker.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(Equal(broker))
			})
		})

		Context("when there is no brokers registered", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should return 404", func() {
				_, err := client.GetBrokerByID(broker.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+broker.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.GetBrokerByID(broker.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+broker.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.GetBrokerByID(broker.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+broker.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

			})
		})
	})

	Describe("Poll", func() {
		Context("when there is operation registered with this id", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(operation)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return it", func() {
				result, err := client.Poll(web.ServiceBrokersURL+"/"+broker.ID+"/"+operation.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(Equal(operation))
			})
		})

		Context("when there is no operation with this id", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should return 404", func() {
				_, err := client.Poll(web.ServiceBrokersURL+"/"+broker.ID+"/"+operation.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), web.ServiceBrokersURL+"/"+broker.ID+"/"+operation.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.Poll(web.ServiceBrokersURL+"/"+broker.ID+"/"+operation.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), web.ServiceBrokersURL+"/"+broker.ID+"/"+operation.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.Poll(web.ServiceBrokersURL+"/"+broker.ID+"/"+operation.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), web.ServiceBrokersURL+"/"+broker.ID+"/"+operation.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

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
				result, err := client.ListPlatforms(params)
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
			It("should return an empty array", func() {
				result, err := client.ListPlatforms(params)
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
				_, err := client.ListPlatforms(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListPlatforms(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("List Visibilities", func() {
		Context("when there are visibilities registered", func() {
			BeforeEach(func() {
				visibilitiesArray := []types.Visibility{*visibility}
				visibilities := types.Visibilities{Visibilities: visibilitiesArray}
				responseBody, _ := json.Marshal(visibilities)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all", func() {
				result, err := client.ListVisibilities(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Visibilities).To(HaveLen(1))
				Expect(result.Visibilities[0]).To(Equal(*visibility))
			})
		})

		Context("when there are no visibilities registered", func() {
			BeforeEach(func() {
				visibilitiesArray := []types.Visibility{}
				visibilities := types.Visibilities{Visibilities: visibilitiesArray}
				responseBody, _ := json.Marshal(visibilities)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return an empty array", func() {
				result, err := client.ListVisibilities(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Visibilities).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListVisibilities(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListVisibilities(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("List service instances", func() {
		Context("when there are service instances registered", func() {
			BeforeEach(func() {
				instancesArray := []types.ServiceInstance{*instance}
				instances := types.ServiceInstances{ServiceInstances: instancesArray}
				responseBody, _ := json.Marshal(instances)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceInstancesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all", func() {
				result, err := client.ListInstances(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.ServiceInstances).To(HaveLen(1))
				Expect(result.ServiceInstances[0]).To(Equal(*instance))
			})
		})

		Context("when there are no service instances registered", func() {
			BeforeEach(func() {
				instancesArray := []types.ServiceInstance{}
				instances := types.ServiceInstances{ServiceInstances: instancesArray}
				responseBody, _ := json.Marshal(instances)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceInstancesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return an empty array", func() {
				result, err := client.ListInstances(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.ServiceInstances).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceInstancesURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListInstances(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceInstancesURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListInstances(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Get service instance", func() {
		Context("when there is instance with this id", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(instance)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceInstancesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return it", func() {
				result, err := client.GetInstanceByID(instance.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(Equal(instance))
			})
		})

		Context("when there is no instance with this id", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceInstancesURL + "/", ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should return 404", func() {
				_, err := client.GetInstanceByID(instance.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+instance.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceInstancesURL + "/", ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.GetInstanceByID(instance.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+instance.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceInstancesURL + "/", ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.GetInstanceByID(instance.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+instance.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

			})
		})
	})

	Describe("List service bindings", func() {
		Context("when there are service bindings registered", func() {
			BeforeEach(func() {
				bindingsArray := []types.ServiceBinding{*binding}
				bindings := types.ServiceBindings{ServiceBindings: bindingsArray}
				responseBody, _ := json.Marshal(bindings)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBindingsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all", func() {
				result, err := client.ListBindings(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.ServiceBindings).To(HaveLen(1))
				Expect(result.ServiceBindings[0]).To(Equal(*binding))
			})
		})

		Context("when there are no service bindings registered", func() {
			BeforeEach(func() {
				bindingsArray := []types.ServiceBinding{}
				bindings := types.ServiceBindings{ServiceBindings: bindingsArray}
				responseBody, _ := json.Marshal(bindings)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBindingsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return an empty array", func() {
				result, err := client.ListBindings(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.ServiceBindings).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBindingsURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListBindings(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBindingsURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListBindings(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Marketplace", func() {
		Context("when there are offerings provided", func() {
			BeforeEach(func() {
				offerings := types.ServiceOfferings{ServiceOfferings: []types.ServiceOffering{*initialOffering}}
				offeringResponseBody, _ := json.Marshal(offerings)

				plans := types.ServicePlans{ServicePlans: []types.ServicePlan{*plan}}
				plansResponseBody, _ := json.Marshal(plans)

				brokerResponseBody, _ := json.Marshal(broker)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseBody: offeringResponseBody, ResponseStatusCode: http.StatusOK},
					{Method: http.MethodGet, Path: web.ServicePlansURL, ResponseBody: plansResponseBody, ResponseStatusCode: http.StatusOK},
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseBody: brokerResponseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all with plans and broker name populated", func() {
				result, err := client.Marketplace(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.ServiceOfferings).To(HaveLen(1))
				Expect(result.ServiceOfferings[0]).To(Equal(*resultOffering))
			})
		})

		Context("when there are no offerings provided", func() {
			BeforeEach(func() {
				offerings := types.ServiceOfferings{}
				offeringResponseBody, _ := json.Marshal(offerings)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseBody: offeringResponseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return an empty array", func() {
				result, err := client.Marketplace(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.ServiceOfferings).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.Marketplace(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.Marketplace(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

	})

	Describe("Delete brokers", func() {
		Context("when an existing broker is being deleted", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteBrokers(params)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteBrokers(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Broker not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteBrokers(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Delete platforms", func() {
		Context("when an existing platform is being deleted", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeletePlatforms(params)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeletePlatforms(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Platform not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeletePlatforms(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Delete visibility", func() {
		Context("when an existing visibility is being deleted", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteVisibilities(params)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteVisibilities(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Visibility not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteVisibilities(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Update brokers", func() {
		Context("when an existing broker is being updated synchronously", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully updated", func() {
				updatedBroker, location, err := client.UpdateBroker("id", broker, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				Expect(updatedBroker).To(Equal(broker))
			})
		})

		Context("when an existing broker is being updated asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "test-location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should be successfully updated", func() {
				updatedBroker, location, err := client.UpdateBroker("id", broker, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
				Expect(updatedBroker).To(BeNil())
			})
		})

		Context("when a non-existing broker is being updated", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"description": "Broker not found"}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				_, location, err := client.UpdateBroker("id", broker, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				_, location, err := client.UpdateBroker("id", broker, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Update platforms", func() {
		Context("when an existing platform is being updated", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(platform)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				updatedPlatform, err := client.UpdatePlatform("id", platform, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updatedPlatform).To(Equal(platform))
			})
		})

		Context("when a non-existing platform is being updated", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"description": "Platform not found"}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdatePlatform("id", platform, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdatePlatform("id", platform, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Update visibilities", func() {
		Context("when an existing visibility is being updated", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(visibility)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				updatedVisibility, err := client.UpdateVisibility("id", visibility, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updatedVisibility).To(Equal(visibility))
			})
		})

		Context("when a non-existing visibility is being updated", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"description": "Visibility not found"}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdateVisibility("id", visibility, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdateVisibility("id", visibility, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

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
				info, _ := client.GetInfo(params)
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
				info, _ := client.GetInfo(params)
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
				_, err := client.GetInfo(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

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
				_, err := client.GetInfo(params)
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Label", func() {
		Context("when valid label change is sent", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusOK}}
			})
			It("should label resource sucessfully", func() {

				err := client.Label(web.ServiceBrokersURL, "id", labelChanges, params)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusNotFound}}
			})
			It("should return error", func() {
				err := client.Label(web.ServiceBrokersURL, "id", labelChanges, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = NewClient(fakeAuthClient, "invalidURL")
				err := client.Label(web.ServiceBrokersURL, "id", labelChanges, params)
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
