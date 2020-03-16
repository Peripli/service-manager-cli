package test

import (
	"encoding/json"
	"github.com/Peripli/service-manager/pkg/web"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Poll test", func() {
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
