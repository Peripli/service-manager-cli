package test

import (
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"net/http"

	"github.com/Peripli/service-manager/pkg/web"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Label test", func() {
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
			client = smclient.NewClient(fakeAuthClient, "invalidURL")
			err := client.Label(web.ServiceBrokersURL, "id", labelChanges, params)
			Expect(err).Should(HaveOccurred())
		})
	})
})
