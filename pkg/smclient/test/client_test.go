package test

import (
	"context"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/Peripli/service-manager/pkg/web"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SM Client test", func() {
	Describe("Test failing client authentication", func() {

		Context("When wrong token is used", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL},
				}
			})
			It("should fail to authentication", func() {
				fakeAuthClient.AccessToken = invalidToken
				client = smclient.NewClient(context.TODO(), fakeAuthClient, smServer.URL)
				_, err := client.ListBrokers(params)

				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, []byte{}, http.StatusUnauthorized)
			})
		})
	})
})
