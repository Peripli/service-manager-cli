package test

import (
	"encoding/json"
	cliquery "github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager/pkg/web"
	"net/http"
	"strings"

	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marketplace test", func() {
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

	Context("when there any supported platform parameter provided", func() {
		of1 := &types.ServiceOffering{
			ID:          "of1",
			Name:        "of1",
			Description: "of1",
			BrokerID:    "of1",
		}
		of2 := &types.ServiceOffering{
			ID:          "of2",
			Name:        "of2",
			Description: "of2",
			BrokerID:    "of2",
		}
		of3 := &types.ServiceOffering{
			ID:          "of3",
			Name:        "of3",
			Description: "of3",
			BrokerID:    "of3",
		}
		of4 := &types.ServiceOffering{
			ID:          "of4",
			Name:        "of4",
			Description: "of4",
			BrokerID:    "of4",
		}
		of5 := &types.ServiceOffering{
			ID:          "of5",
			Name:        "of5",
			Description: "of5",
			BrokerID:    "of5",
		}
		plan_supports_any_env := &types.ServicePlan{
			ID:                "plan_supports_any_env",
			Name:              "plan_supports_any_env",
			Description:       "no metadata plan", // supports any env
			ServiceOfferingID: "of1",
		}
		plan_supports_cf_env1 := &types.ServicePlan{
			ID:                "plan_supports_cf_env1",
			Name:              "plan_supports_cf_env1",
			Description:       "plan with cf support",
			ServiceOfferingID: "of2",
			Metadata:          []byte("{\"supportedPlatforms\":[\"cloudfoundry\"]}"),
		}
		plan_supports_cf_env2 := &types.ServicePlan{
			ID:                "plan_supports_cf_env2",
			Name:              "plan_supports_cf_env2",
			Description:       "plan with cf support",
			ServiceOfferingID: "of3",
			Metadata:          []byte("{\"supportedPlatforms\":[\"cloudfoundry\"]}"),
		}
		plan_supports_k8s := &types.ServicePlan{
			ID:                "plan_supports_k8s",
			Name:              "plan_supports_k8s",
			Description:       "plan with k8s support",
			ServiceOfferingID: "of3",
			Metadata:          []byte("{\"supportedPlatforms\":[\"kubernetes\"]}"),
		}
		plan_supports_k8s_and_cf := &types.ServicePlan{
			ID:                "plan_supports_k8s_and_cf",
			Name:              "plan_supports_k8s_and_cf",
			Description:       "plan with cf & k8s support",
			ServiceOfferingID: "of4",
			Metadata:          []byte("{\"supportedPlatforms\":[\"kubernetes\",\"cloudfoundry\"]}"),
		}
		plan_supports_test_env := &types.ServicePlan{
			ID:                "plan_supports_test_env",
			Name:              "plan_supports_test_env",
			Description:       "plan with test env support",
			ServiceOfferingID: "of5",
			Metadata:          []byte("{\"supportedPlatforms\":[\"test-env\"]}"),
		}
		rof1_with_plan_any_env := &types.ServiceOffering{
			ID:          "of1",
			Name:        "of1",
			Description: "of1",
			BrokerID:    "of1",
			Plans:       []types.ServicePlan{*plan_supports_any_env},
		}
		rof2_with_plan_cf_env := &types.ServiceOffering{
			ID:          "of2",
			Name:        "of2",
			Description: "of2",
			BrokerID:    "of2",
			Plans:       []types.ServicePlan{*plan_supports_cf_env1},
		}
		rof3_with_plan_cf_and_k8s := &types.ServiceOffering{
			ID:          "of3",
			Name:        "of3",
			Description: "of3",
			BrokerID:    "of3",
			Plans:       []types.ServicePlan{*plan_supports_cf_env2, *plan_supports_k8s},
		}
		rof4_with_plan_cf_and_k8s := &types.ServiceOffering{
			ID:          "of4",
			Name:        "of4",
			Description: "of4",
			BrokerID:    "of4",
			Plans:       []types.ServicePlan{*plan_supports_k8s_and_cf},
		}
		rof5_with_plan_test_env := &types.ServiceOffering{
			ID:          "of5",
			Name:        "of5",
			Description: "of5",
			BrokerID:    "of5",
			Plans:       []types.ServicePlan{*plan_supports_test_env},
		}

		BeforeEach(func() {

			offerings := types.ServiceOfferings{ServiceOfferings: []types.ServiceOffering{
				*of1,
				*of2,
				*of3,
				*of4,
				*of5,
			}}
			offeringResponseBody, _ := json.Marshal(offerings)

			brokerResponseBody, _ := json.Marshal(broker)

			handlerDetails = []HandlerDetails{
				{Method: http.MethodGet, Path: web.ServiceOfferingsURL, ResponseBody: offeringResponseBody, ResponseStatusCode: http.StatusOK},
				{
					Method:             http.MethodGet,
					Path:               web.ServicePlansURL,
					ResponseStatusCode: http.StatusOK,
					ResponseBodyProvider: func(req *http.Request) []byte {
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of1%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*plan_supports_any_env},
							})
							return response
						}
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of2%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*plan_supports_cf_env1},
							})
							return response
						}
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of3%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*plan_supports_cf_env2, *plan_supports_k8s},
							})
							return response
						}
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of4%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*plan_supports_k8s_and_cf},
							})
							return response
						}
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of5%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*plan_supports_test_env},
							})
							return response
						}
						return []byte{}
					}},
				{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseBody: brokerResponseBody, ResponseStatusCode: http.StatusOK},
			}
		})
		It("should return all with all plans if any environment requested", func() {
			result, err := client.Marketplace(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "any",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServiceOfferings).To(HaveLen(5))
			Expect(result.ServiceOfferings[0]).To(Equal(*rof1_with_plan_any_env))
			Expect(result.ServiceOfferings[1]).To(Equal(*rof2_with_plan_cf_env))
			Expect(result.ServiceOfferings[2]).To(Equal(*rof3_with_plan_cf_and_k8s))
			Expect(result.ServiceOfferings[3]).To(Equal(*rof4_with_plan_cf_and_k8s))
			Expect(result.ServiceOfferings[4]).To(Equal(*rof5_with_plan_test_env))
		})
		It("should return cloudfoundry services with all plans if cloudfoundry environment requested", func() {
			result, err := client.Marketplace(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "cloudfoundry",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServiceOfferings).To(HaveLen(4))
			Expect(result.ServiceOfferings[0]).To(Equal(*rof1_with_plan_any_env)) // matches as it any
			Expect(result.ServiceOfferings[1]).To(Equal(*rof2_with_plan_cf_env))
			Expect(result.ServiceOfferings[2]).To(Equal(*rof3_with_plan_cf_and_k8s))
			Expect(result.ServiceOfferings[3]).To(Equal(*rof4_with_plan_cf_and_k8s))
		})
		It("should return kubernetes services with all plans if kubernetes environment requested", func() {
			result, err := client.Marketplace(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "kubernetes",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServiceOfferings).To(HaveLen(3))
			Expect(result.ServiceOfferings[0]).To(Equal(*rof1_with_plan_any_env)) // matches as it any
			Expect(result.ServiceOfferings[1]).To(Equal(*rof3_with_plan_cf_and_k8s))
			Expect(result.ServiceOfferings[2]).To(Equal(*rof4_with_plan_cf_and_k8s))
		})
		It("should return test-env services with all plans if test-env environment requested", func() {
			result, err := client.Marketplace(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "test-env",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServiceOfferings).To(HaveLen(2))
			Expect(result.ServiceOfferings[0]).To(Equal(*rof1_with_plan_any_env)) // matches as it any
			Expect(result.ServiceOfferings[1]).To(Equal(*rof5_with_plan_test_env))
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
