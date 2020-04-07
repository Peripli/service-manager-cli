package test

import (
	"context"
	"encoding/json"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"net/http"

	"github.com/Peripli/service-manager/pkg/web"

	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Platform test", func() {
	platform := &types.Platform{
		ID:          "platformID",
		Name:        "cfeu10",
		Type:        "cf",
		Description: "Test platform",
	}

	Describe("List platforms", func() {
		Context("when there are platforms registered", func() {
			BeforeEach(func() {
				platformsArray := []types.Platform{*platform}
				platforms := types.Platforms{Platforms: platformsArray}
				responseBody, _ := json.Marshal(platforms)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return all", func() {
				result, err := client.ListPlatforms(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Platforms).To(HaveLen(1))
				Expect(result.Platforms[0]).To(Equal(*platform))
			})
		})

		Context("when there are no platforms registered", func() {
			BeforeEach(func() {
				platformsArray := []types.Platform{}
				platforms := types.Platforms{Platforms: platformsArray}
				responseBody, _ := json.Marshal(platforms)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should return an empty array", func() {
				result, err := client.ListPlatforms(params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Platforms).To(HaveLen(0))
			})
		})

		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle status code != 200", func() {
				_, err := client.ListPlatforms(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
		Context("when invalid status code is returned", func() {
			BeforeEach(func() {
				handlerDetails = []HandlerDetails{
					{Method: http.MethodGet, Path: web.PlatformsURL, ResponseStatusCode: http.StatusBadRequest},
				}
			})
			It("should handle status code > 299", func() {
				_, err := client.ListPlatforms(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Register platform", func() {
		Context("When valid platform is being registered", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(platform)
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should register successfully", func() {
				responsePlatform, err := client.RegisterPlatform(platform, params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responsePlatform).To(Equal(platform))
			})
		})

		Context("When invalid platform is returned by SM", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(struct {
					Name bool `json:"name"`
				}{
					Name: true,
				})
				handlerDetails = []HandlerDetails{
					{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should return error", func() {
				responsePlatform, err := client.RegisterPlatform(platform, params)

				Expect(err).Should(HaveOccurred())
				Expect(responsePlatform).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SM", func() {
			Context("And status code is successful", func() {
				BeforeEach(func() {
					responseBody, _ := json.Marshal(platform)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
					}
				})
				It("should return error with status code", func() {
					responsePlatform, err := client.RegisterPlatform(platform, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responsePlatform).To(BeNil())
				})
			})

			Context("And status code is unsuccessful", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error with url and description", func() {
					responsePlatform, err := client.RegisterPlatform(platform, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responsePlatform).To(BeNil())
				})
			})

			Context("And response body is invalid", func() {
				BeforeEach(func() {
					responseBody := []byte(`{ "description": description", "error": "error"}`)
					handlerDetails = []HandlerDetails{
						{Method: http.MethodPost, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusBadRequest},
					}
				})
				It("should return error without url and description if invalid response body", func() {
					responsePlatform, err := client.RegisterPlatform(platform, params)

					Expect(err).Should(HaveOccurred())
					verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
					Expect(responsePlatform).To(BeNil())
				})
			})
		})

		Context("When invalid config is set", func() {
			It("should return error", func() {
				client = smclient.NewClient(context.TODO(), fakeAuthClient, "invalidURL")
				_, err := client.RegisterPlatform(platform, params)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Update platforms", func() {
		Context("when an existing platform is being updated", func() {
			BeforeEach(func() {
				responseBody, _ := json.Marshal(platform)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				updatedPlatform, err := client.UpdatePlatform("id", platform, params)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updatedPlatform).To(Equal(platform))
			})
		})

		Context("when a non-existing platform is being updated", func() {
			BeforeEach(func() {
				responseBody := []byte(`{"description": "Platform not found"}`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdatePlatform("id", platform, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodPatch, Path: web.PlatformsURL + "/", ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				_, err := client.UpdatePlatform("id", platform, params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path+"id", handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})

	Describe("Delete platforms", func() {
		Context("when an existing platform is being deleted", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusOK},
				}
			})
			It("should be successfully removed", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeletePlatforms(params)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when service manager returns a non-expected status code", func() {
			BeforeEach(func() {
				responseBody := []byte("{}")

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusCreated},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeletePlatforms(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})

		Context("when service manager returns a status code not found", func() {
			BeforeEach(func() {
				responseBody := []byte(`{ "description": "Platform not found" }`)

				handlerDetails = []HandlerDetails{
					{Method: http.MethodDelete, Path: web.PlatformsURL, ResponseBody: responseBody, ResponseStatusCode: http.StatusNotFound},
				}
			})
			It("should handle error", func() {
				params.FieldQuery = append(params.FieldQuery, "id eq 'id'")
				err := client.DeletePlatforms(params)
				Expect(err).Should(HaveOccurred())
				verifyErrorMsg(err.Error(), handlerDetails[0].Path, handlerDetails[0].ResponseBody, handlerDetails[0].ResponseStatusCode)
			})
		})
	})
})
