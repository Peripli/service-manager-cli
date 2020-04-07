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

var _ = Describe("Broker test", func() {
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
				result, err := client.ListBrokers(params)
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
			It("should return an empty array", func() {
				result, err := client.ListBrokers(params)
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
				_, err := client.ListBrokers(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListBrokers(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

			})
		})
	})

	Describe("Get broker", func() {
		Context("when there is broker registered with this id", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return it", func() {
				result, err := client.GetBrokerByID(broker.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(Equal(broker))
			})
		})

		Context("when there is no brokers registered", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should return 404", func() {
				_, err := client.GetBrokerByID(broker.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+broker.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.GetBrokerByID(broker.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+broker.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.GetBrokerByID(broker.ID, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+broker.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

			})
		})
	})

	Describe("Register broker", func() {
		Context("When valid broker is being registered synchronously", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responseBroker, location, err := client.RegisterBroker(broker, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				Expect(responseBroker).To(Equal(broker))
			})
		})

		Context("When valid broker is being registered asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "test-location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.ServiceBrokersURL, ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should receive operation location", func() {
				responseBroker, location, err := client.RegisterBroker(broker, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
				Expect(responseBroker).To(BeNil())
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
				responseBroker, location, err := client.RegisterBroker(broker, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
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
					responseBroker, location, err := client.RegisterBroker(broker, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
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
					responseBroker, location, err := client.RegisterBroker(broker, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
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
					responseBroker, location, err := client.RegisterBroker(broker, params)

					Expect(err).Should(HaveOccurred())
					Expect(location).Should(BeEmpty())

					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseBroker).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = smclient.NewClient(context.TODO(),fakeAuthClient, "invalidURL")
				_, location, err := client.RegisterBroker(broker, params)

				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
			})
		})
	})

	Describe("Update brokers", func() {
		Context("when an existing broker is being updated synchronously", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(broker)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully updated", func() {
				updatedBroker, location, err := client.UpdateBroker("id", broker, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				Expect(updatedBroker).To(Equal(broker))
			})
		})

		Context("when an existing broker is being updated asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "test-location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should be successfully updated", func() {
				updatedBroker, location, err := client.UpdateBroker("id", broker, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
				Expect(updatedBroker).To(BeNil())
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
				_, location, err := client.UpdateBroker("id", broker, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
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
				_, location, err := client.UpdateBroker("id", broker, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(HaveLen(0))
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Delete broker", func() {
		Context("when an existing broker is being deleted synchronously", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBrokersURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				location, err := client.DeleteBroker(broker.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(BeEmpty())
			})
		})

		Context("when an existing broker is being deleted asynchronously", func() {
			var locationHeader string
			BeforeEach(func() {
				locationHeader = "location"
				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.ServiceBrokersURL + "/", ResponseStatusCode: http.StatusAccepted, Headers: map[string]string{"Location": locationHeader}},
				}
			})
			It("should be successfully removed", func() {
				location, err := client.DeleteBroker(broker.ID, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(location).Should(Equal(locationHeader))
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
				location, err := client.DeleteBroker(broker.ID, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+broker.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
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
				location, err := client.DeleteBroker(broker.ID, params)
				Expect(err).Should(HaveOccurred())
				Expect(location).Should(BeEmpty())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+broker.ID, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})
})
