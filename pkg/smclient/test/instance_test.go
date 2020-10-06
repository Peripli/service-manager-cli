package test

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/smclient"

	"github.com/Peripli/service-manager/pkg/web"

	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Instance test", func() {
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

	Describe("Get service instance parameters", func() {
		When("there is instance with this id with parameters", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(instanceParameters)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet,
						Path: web.ServiceInstancesURL + "/"+ instance.ID + web.ParametersURL,
						ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return parameters", func() {
				result, err := client.GetInstanceParameters(instance.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(Equal(instanceParameters))
			})
		})

		When("there is instance with this id without parameters", func() {
			instanceParameters := make(map[string]interface{})
			BeforeEach(func() {
				responseBody, _ := json.Marshal(instanceParameters)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet,
						Path: web.ServiceInstancesURL + "/"+ instance.ID + web.ParametersURL,
						ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return empty parameters", func() {
				result, err := client.GetInstanceParameters(instance.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(Equal(instanceParameters))
			})
		})

		When("there is no instance with this id", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet,
						Path: web.ServiceInstancesURL + "/"+ instance.ID + web.ParametersURL,
						ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should return 404", func() {
				_, err := client.GetInstanceParameters(instance.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		When("bad gateway status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet,
						Path: web.ServiceInstancesURL + "/"+ instance.ID + web.ParametersURL,
						ResponseStatusCode: http.StatusBadGateway},
				}
			})
			It("should return an error with status code 502", func() {
				_, err := client.GetInstanceParameters(instance.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		When("bad request status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet,
						Path: web.ServiceInstancesURL + "/"+ instance.ID + web.ParametersURL,
						ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should return an error with status code 400", func() {
				_, err := client.GetInstanceParameters(instance.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

			})
		})
	})

	Describe("Provision", func() {
		Context("When valid instance is being provisioned synchronously", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(instance)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceInstancesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should provision successfully", func() {
				responseInstance, location, err := client.Provision(instance, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				Expect(responseInstance).To(Equal(instance))
			})
		})

		Context("When valid instance is being provisioned asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "test-location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceInstancesURL, ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should receive operation location", func() {
				responseInstance, location, err := client.Provision(instance, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
				Expect(responseInstance).To(BeNil())
			})
		})

		Context("When invalid instance is being returned by SM", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(struct {
					Name bool `json:"name"`
				}{
					Name: true,
				})
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceInstancesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should return error", func() {
				responseInstance, location, err := client.Provision(instance, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				Expect(responseInstance).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(instance)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceInstancesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
					}
				})
				It("should return error with status code", func() {
					responseInstance, location, err := client.Provision(instance, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseInstance).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceInstancesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responseInstance, location, err := client.Provision(instance, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseInstance).To(BeNil())
				})
			})

			Context("And invalid response body", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.ServiceInstancesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responseInstance, location, err := client.Provision(instance, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())

					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseInstance).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = smclient.NewClient(context.TODO(), fakeAuthClient, "invalidURL")
				_, location, err := client.Provision(instance, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
			})
		})
	})

	Describe("Deprovision", func() {
		Context("when an existing instance is being deleted synchronously", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")
				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceInstancesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				location, err := client.Deprovision(instance.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(BeEmpty())
			})
		})

		Context("when an existing instance is being deleted asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceInstancesURL + "/", ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should be successfully removed", func() {
				location, err := client.Deprovision(instance.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceInstancesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				location, err := client.Deprovision(instance.ID, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+instance.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Broker not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceInstancesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				location, err := client.Deprovision(instance.ID, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+instance.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Update", func() {
		Context("When valid instance is being updated synchronously", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(instance)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceInstancesURL + "/" + instance.ID, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should update successfully", func() {
				responseInstance, location, err := client.UpdateInstance(instance.ID, instance, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				Expect(responseInstance).To(Equal(instance))
			})
		})

		Context("When valid instance is being updated asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "test-location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceInstancesURL + "/" + instance.ID, ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should receive operation location", func() {
				responseInstance, location, err := client.UpdateInstance(instance.ID, instance, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
				Expect(responseInstance).To(BeNil())
			})
		})

		Context("When invalid instance is being returned by SM", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(struct {
					Name bool `json:"name"`
				}{
					Name: true,
				})
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceInstancesURL + "/" + instance.ID, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return error", func() {
				responseInstance, location, err := client.UpdateInstance(instance.ID, instance, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				Expect(responseInstance).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(instance)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPatch, Path: web.ServiceInstancesURL + "/" + instance.ID, ResponseBody: responseBody, ResponseStatusCode: http.StatusTeapot},
					}
				})
				It("should return error with status code", func() {
					responseInstance, location, err := client.UpdateInstance(instance.ID, instance, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseInstance).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPatch, Path: web.ServiceInstancesURL + "/" + instance.ID, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responseInstance, location, err := client.UpdateInstance(instance.ID, instance, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseInstance).To(BeNil())
				})
			})

			Context("And invalid response body", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPatch, Path: web.ServiceInstancesURL + "/" + instance.ID, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responseInstance, location, err := client.UpdateInstance(instance.ID, instance, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())

					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseInstance).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = smclient.NewClient(context.TODO(), fakeAuthClient, "invalidURL")
				_, location, err := client.UpdateInstance(instance.ID, instance, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
			})
		})
	})
})
