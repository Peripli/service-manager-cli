package test

import (
	"context"
	"encoding/json"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/Peripli/service-manager/pkg/web"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Binding test", func() {
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

	Describe("Get service binding", func() {
		Context("when there is binding with this id", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(binding)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBindingsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return it", func() {
				result, err := client.GetBindingByID(binding.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(Equal(binding))
			})
		})

		Context("when there is no binding with this id", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBindingsURL + "/", ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should return 404", func() {
				_, err := client.GetBindingByID(binding.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+binding.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBindingsURL + "/", ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.GetBindingByID(binding.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+binding.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBindingsURL + "/", ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.GetBindingByID(binding.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+binding.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

			})
		})
	})

	Describe("Bind", func() {
		Context("When valid binding is being created synchronously", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(binding)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceBindingsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should provision successfully", func() {
				responseBinding, location, err := client.Bind(binding, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				Expect(responseBinding).To(Equal(binding))
			})
		})

		Context("When valid binding is being created asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "test-location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceBindingsURL, ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should receive operation location", func() {
				responseBinding, location, err := client.Bind(binding, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
				Expect(responseBinding).To(BeNil())
			})
		})

		Context("When invalid binding is being returned by SM", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(struct {
					Name bool `json:"name"`
				}{
					Name: true,
				})
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceBindingsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should return error", func() {
				responseBinding, location, err := client.Bind(binding, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				Expect(responseBinding).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(binding)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceBindingsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
					}
				})
				It("should return error with status code", func() {
					responseBinding, location, err := client.Bind(binding, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseBinding).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceBindingsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responseBinding, location, err := client.Bind(binding, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseBinding).To(BeNil())
				})
			})

			Context("And invalid response body", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceBindingsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responseBinding, location, err := client.Bind(binding, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())

					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseBinding).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = smclient.NewClient(context.TODO(), fakeAuthClient, "invalidURL")
				_, location, err := client.Bind(binding, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
			})
		})
	})

	Describe("Unbind", func() {
		Context("when an existing binding is being deleted synchronously", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")
				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBindingsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				location, err := client.Unbind(binding.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(BeEmpty())
			})
		})

		Context("when an existing binding is being deleted asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBindingsURL + "/", ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should be successfully removed", func() {
				location, err := client.Unbind(binding.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBindingsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				location, err := client.Unbind(binding.ID, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+binding.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Broker not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBindingsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				location, err := client.Unbind(binding.ID, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+binding.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})
})
