package smclient

import (
	"encoding/json"
	"fmt"
	cliquery "github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager/pkg/query"
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
	var fakeAuthClient *FakeAuthClient

	platform := &types.Platform{
		ID:          "platformID",
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

	labelChanges := &types.LabelChanges{
		LabelChanges: []*query.LabelChange{
			{Key: "key", Operation: query.LabelOperation("add"), Values: []string{"val1", "val2"}},
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
		fakeAuthClient = &FakeAuthClient{AccessToken: validToken}
		client = NewClient(fakeAuthClient, smServer.URL)
	})

	Describe("ParseQuery encodes", func() {

		input := [][]string{{"description = description with multiple     spaces"},
			{"description = description with operators: [in = != eqornil gt lt in notin]"},
			{`description = description with \`},
			{`description in [description with "quotes"||description with \]`},
			{"type = type", `description = description with "quotes"`}}
		output := []string{"description+%3D+description+with+multiple+++++spaces",
			"description+%3D+description+with+operators%3A+%5Bin+%3D+%21%3D+eqornil+gt+lt+in+notin%5D",
			"description+%3D+description+with+%5C",
			"description+in+%5Bdescription+with+%22quotes%22%7C%7Cdescription+with+%5C%5D",
			"type+%3D+type|description+%3D+description+with+%22quotes%22"}

		Context("when queries are provided", func() {
			It("should url encode and join them", func() {
				for i := range input {
					Expect(parseQuery(input[i])).To(Equal(output[i]))
				}
			})
		})
	})

	Describe("General parameter", func() {
		Context("In query parameters", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"token_issuer_url": "http://uaa.com"}`)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.InfoURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should make request with these parameters", func() {
				param := cliquery.Parameters{}
				param.Add(cliquery.GeneralParameter, "key=val")
				info, err := client.GetInfo(param.Copy())
				Expect(err).ToNot(HaveOccurred())
				Expect(info).ToNot(BeNil())
				Expect(fakeAuthClient.requestURI).To(Equal(fmt.Sprintf("%s?key=val", web.InfoURL)))
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
				client = NewClient(http.DefaultClient, smServer.URL)
				_, err := client.ListBrokers()

				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.ServiceBrokersURL, StatusCode: http.StatusUnauthorized}))
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

	Describe("Register Visibility", func() {
		Context("When valid visibility is being registered", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(visibility)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responseVisibility, err := client.RegisterVisibility(visibility)

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
				responseVisibility, err := client.RegisterVisibility(visibility)

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
					responseVisibility, err := client.RegisterVisibility(visibility)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{StatusCode: handlerDetails[0].ResponseStatusCode}))
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
					responseVisibility, err := client.RegisterVisibility(visibility)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.VisibilitiesURL, Description: "error", StatusCode: handlerDetails[0].ResponseStatusCode}))
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
					responseVisibility, err := client.RegisterVisibility(visibility)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.VisibilitiesURL, StatusCode: handlerDetails[0].ResponseStatusCode}))
					Expect(responseVisibility).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = NewClient(http.DefaultClient, "invalidURL")
				_, err := client.RegisterVisibility(visibility)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Register broker", func() {
		Context("When valid broker is being registered", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
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
					{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
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
						{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
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
						{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responseBroker, err := client.RegisterBroker(broker)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.ServiceBrokersURL, Description: "description",
						ErrorMessage: "error", StatusCode: handlerDetails[0].ResponseStatusCode}))
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
					responseBroker, err := client.RegisterBroker(broker)

					Expect(err).Should(HaveOccurred())
					Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.ServiceBrokersURL, StatusCode: handlerDetails[0].ResponseStatusCode}))
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
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
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
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
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
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseStatusCode: http.StatusCreated},
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
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListBrokers()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusBadRequest, URL: smServer.URL + web.ServiceBrokersURL}))
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
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusBadRequest, URL: smServer.URL + web.PlatformsURL}))
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
				result, err := client.ListVisibilities()
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
			It("should return empty array", func() {
				result, err := client.ListVisibilities()
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
				_, err := client.ListVisibilities()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListVisibilities()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{StatusCode: http.StatusBadRequest, URL: smServer.URL + web.VisibilitiesURL}))
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

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseBody: offeringResponseBody, ResponseStatusCode: http.StatusOK},
					{Method: http.MethodGet, Path: web.ServicePlansURL, ResponseBody: plansResponseBody, ResponseStatusCode: http.StatusOK},
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseBody: brokerResponseBody, ResponseStatusCode: http.StatusOK},
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

				handlerDetails = []HandlerDetails{
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
				handlerDetails = []HandlerDetails{
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
				handlerDetails = []HandlerDetails{
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
					{Method: http.MethodDelete, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
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
					{Method: http.MethodDelete, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
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
					{Method: http.MethodDelete, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				err := client.DeleteBroker("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Broker not found", URL: smServer.URL + web.ServiceBrokersURL + "/id", StatusCode: http.StatusNotFound}))
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

	Describe("Delete visibility", func() {
		Context("when an existing visibility is being deleted", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				err := client.DeleteVisibility("id")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				err := client.DeleteVisibility("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Visibility not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				err := client.DeleteVisibility("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Visibility not found", URL: smServer.URL + web.VisibilitiesURL + "/id", StatusCode: http.StatusNotFound}))
			})
		})
	})

	Describe("Update brokers", func() {
		Context("when an existing broker is being updated", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				updatedBroker, err := client.UpdateBroker("id", broker)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updatedBroker).To(Equal(broker))
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
				_, err := client.UpdateBroker("id", broker)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Broker not found", URL: smServer.URL + web.ServiceBrokersURL + "/id", StatusCode: http.StatusNotFound}))
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
				_, err := client.UpdateBroker("id", broker)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
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
				updatedPlatform, err := client.UpdatePlatform("id", platform)
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
				_, err := client.UpdatePlatform("id", platform)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Platform not found", URL: smServer.URL + web.PlatformsURL + "/id", StatusCode: http.StatusNotFound}))
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
				_, err := client.UpdatePlatform("id", platform)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
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
				updatedVisibility, err := client.UpdateVisibility("id", visibility)
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
				_, err := client.UpdateVisibility("id", visibility)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{Description: "Visibility not found", URL: smServer.URL + web.VisibilitiesURL + "/id", StatusCode: http.StatusNotFound}))
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
				_, err := client.UpdateVisibility("id", visibility)
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(errors.ResponseError{StatusCode: http.StatusCreated}))
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
				info, _ := client.GetInfo(nil)
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
				info, _ := client.GetInfo(nil)
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
				_, err := client.GetInfo(nil)
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
				_, err := client.GetInfo(nil)
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

				err := client.Label(web.ServiceBrokersURL, "id", labelChanges)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusNotFound}}
			})
			It("should return error", func() {
				err := client.Label(web.ServiceBrokersURL, "id", labelChanges)
				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(errors.ResponseError{URL: smServer.URL + web.ServiceBrokersURL + "/id", StatusCode: http.StatusNotFound}))
			})
		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = NewClient(http.DefaultClient, "invalidURL")
				err := client.Label(web.ServiceBrokersURL, "id", labelChanges)
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
