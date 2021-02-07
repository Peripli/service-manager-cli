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
		p1 := &types.ServicePlan{
			ID:                "p1",
			Name:              "p1",
			Description:       "no metadata plan",
			ServiceOfferingID: "of1",
		}
		p2 := &types.ServicePlan{
			ID:                "p2",
			Name:              "p2",
			Description:       "plan with cf support",
			ServiceOfferingID: "of2",
			Metadata:          []byte("{\"supportedPlatforms\":[\"cloudfoundry\"]}"),
		}
		p3 := &types.ServicePlan{
			ID:                "p3",
			Name:              "p3",
			Description:       "plan with cf support",
			ServiceOfferingID: "of3",
			Metadata:          []byte("{\"supportedPlatforms\":[\"cloudfoundry\"]}"),
		}
		p4 := &types.ServicePlan{
			ID:                "p4",
			Name:              "p4",
			Description:       "plan with k8s support",
			ServiceOfferingID: "of3",
			Metadata:          []byte("{\"supportedPlatforms\":[\"kubernetes\"]}"),
		}
		p5 := &types.ServicePlan{
			ID:                "p5",
			Name:              "p5",
			Description:       "plan with cf & k8s support",
			ServiceOfferingID: "of4",
			Metadata:          []byte("{\"supportedPlatforms\":[\"kubernetes\",\"cloudfoundry\"]}"),
		}
		p6 := &types.ServicePlan{
			ID:                "p6",
			Name:              "p6",
			Description:       "plan with test env support",
			ServiceOfferingID: "of5",
			Metadata:          []byte("{\"supportedPlatforms\":[\"test-env\"]}"),
		}
		rof1 := &types.ServiceOffering{
			ID:          "of1",
			Name:        "of1",
			Description: "of1",
			BrokerID:    "of1",
			Plans:       []types.ServicePlan{*p1},
		}
		rof2 := &types.ServiceOffering{
			ID:          "of2",
			Name:        "of2",
			Description: "of2",
			BrokerID:    "of2",
			Plans:       []types.ServicePlan{*p2},
		}
		rof3 := &types.ServiceOffering{
			ID:          "of3",
			Name:        "of3",
			Description: "of3",
			BrokerID:    "of3",
			Plans:       []types.ServicePlan{*p3, *p4},
		}
		rof4 := &types.ServiceOffering{
			ID:          "of4",
			Name:        "of4",
			Description: "of4",
			BrokerID:    "of4",
			Plans:       []types.ServicePlan{*p5},
		}
		rof5 := &types.ServiceOffering{
			ID:          "of5",
			Name:        "of5",
			Description: "of5",
			BrokerID:    "of5",
			Plans:       []types.ServicePlan{*p6},
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
								ServicePlans: []types.ServicePlan{*p1},
							})
							return response
						}
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of2%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*p2},
							})
							return response
						}
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of3%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*p3, *p4},
							})
							return response
						}
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of4%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*p5},
							})
							return response
						}
						if strings.Contains(req.RequestURI, "service_offering_id+eq+%27of5%27") {
							response, _ := json.Marshal(types.ServicePlans{
								ServicePlans: []types.ServicePlan{*p6},
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
			Expect(result.ServiceOfferings[0]).To(Equal(*rof1))
			Expect(result.ServiceOfferings[1]).To(Equal(*rof2))
			Expect(result.ServiceOfferings[2]).To(Equal(*rof3))
			Expect(result.ServiceOfferings[3]).To(Equal(*rof4))
			Expect(result.ServiceOfferings[4]).To(Equal(*rof5))
		})
		It("should return cloudfoundry services with all plans if cloudfoundry environment requested", func() {
			result, err := client.Marketplace(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "cloudfoundry",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServiceOfferings).To(HaveLen(4))
			Expect(result.ServiceOfferings[0]).To(Equal(*rof1)) // matches as it any
			Expect(result.ServiceOfferings[1]).To(Equal(*rof2))
			Expect(result.ServiceOfferings[2]).To(Equal(*rof3))
			Expect(result.ServiceOfferings[3]).To(Equal(*rof4))
		})
		It("should return kubernetes services with all plans if kubernetes environment requested", func() {
			result, err := client.Marketplace(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "kubernetes",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServiceOfferings).To(HaveLen(3))
			Expect(result.ServiceOfferings[0]).To(Equal(*rof1)) // matches as it any
			Expect(result.ServiceOfferings[1]).To(Equal(*rof3))
			Expect(result.ServiceOfferings[2]).To(Equal(*rof4))
		})
		It("should return test-env services with all plans if test-env environment requested", func() {
			result, err := client.Marketplace(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "test-env",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServiceOfferings).To(HaveLen(2))
			Expect(result.ServiceOfferings[0]).To(Equal(*rof1)) // matches as it any
			Expect(result.ServiceOfferings[1]).To(Equal(*rof5))
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
