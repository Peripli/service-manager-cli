package test

import (
	"encoding/json"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/Peripli/service-manager/pkg/web"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Visibility test", func() {

	visibility := &types.Visibility{
		ID:            "visibilityID",
		PlatformID:    "platformID",
		ServicePlanID: "planID",
	}

	Describe("List Visibilities", func() {
		Context("when there are visibilities registered", func() {
			BeforeEach(func() {
				visibilitiesArray := []types.Visibility{*visibility}
				visibilities := types.Visibilities{Visibilities: visibilitiesArray}
				responseBody, _ := json.Marshal(visibilities)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all", func() {
				result, err := client.ListVisibilities(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Visibilities).To(HaveLen(1))
				Expect(result.Visibilities[0]).To(Equal(*visibility))
			})
		})

		Context("when there are no visibilities registered", func() {
			BeforeEach(func() {
				visibilitiesArray := []types.Visibility{}
				visibilities := types.Visibilities{Visibilities: visibilitiesArray}
				responseBody, _ := json.Marshal(visibilities)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return an empty array", func() {
				result, err := client.ListVisibilities(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Visibilities).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListVisibilities(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.VisibilitiesURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListVisibilities(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Register Visibility", func() {
		Context("When valid visibility is being registered", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(visibility)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responseVisibility, err := client.RegisterVisibility(visibility, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responseVisibility).To(Equal(visibility))
			})
		})

		Context("When invalid visibility is returned by SM", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(struct {
					ID bool `json:"id"`
				}{
					ID: true,
				})
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})

			It("should return error", func() {
				responseVisibility, err := client.RegisterVisibility(visibility, params)

				Expect(err).Should(HaveOccurred())
				Expect(responseVisibility).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is successful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(visibility)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
					}
				})
				It("should return error with status code", func() {
					responseVisibility, err := client.RegisterVisibility(visibility, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseVisibility).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responseVisibility, err := client.RegisterVisibility(visibility, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseVisibility).To(BeNil())
				})
			})

			Context("And response body is invalid", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responseVisibility, err := client.RegisterVisibility(visibility, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responseVisibility).To(BeNil())
				})
			})

		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = smclient.NewClient(fakeAuthClient, "invalidURL")
				_, err := client.RegisterVisibility(visibility, params)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Update visibilities", func() {
		Context("when an existing visibility is being updated", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(visibility)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				updatedVisibility, err := client.UpdateVisibility("id", visibility, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updatedVisibility).To(Equal(visibility))
			})
		})

		Context("when a non-existing visibility is being updated", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"description": "Visibility not found"}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdateVisibility("id", visibility, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.VisibilitiesURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdateVisibility("id", visibility, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)

			})
		})
	})

	Describe("Delete visibility", func() {
		Context("when an existing visibility is being deleted", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteVisibilities(params)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteVisibilities(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Visibility not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.VisibilitiesURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeleteVisibilities(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})
})
