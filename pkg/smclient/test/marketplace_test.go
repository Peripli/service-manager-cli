package test

import (
	"encoding/json"
	"github.com/Peripli/service-manager/pkg/web"
	"net/http"

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
