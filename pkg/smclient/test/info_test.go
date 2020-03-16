package test

import (
	"github.com/Peripli/service-manager/pkg/web"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info test", func() {
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

	Context("with general parameter", func() {
		BeforeEach(func() {
			responseBody := []byte(`{"token_issuer_url": "http://uaa.com"}`)
			handlerDetails = []HandlerDetails{
				{Method: http.MethodGet, Path: web.InfoURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
			}
		})

		It("should make request with these parameters", func() {
			info, err := client.GetInfo(params)
			Expect(err).ToNot(HaveOccurred())
			Expect(info).ToNot(BeNil())
		})
	})
})
