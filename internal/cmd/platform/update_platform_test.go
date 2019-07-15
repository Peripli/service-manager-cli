package platform

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"testing"

	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

func TestUpdatePlatformCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Update platform command test", func() {

	var client *smclientfakes.FakeClient
	var command *UpdatePlatformCmd
	var buffer *bytes.Buffer
	var platform types.Platform

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewUpdatePlatformCmd(context)
	})

	validUpdatePlatformExecution := func(args ...string) error {
		platform = types.Platform{
			Name: "platform",
			ID:   "id",
			Type: "type",
		}
		platforms := &types.Platforms{Platforms: []types.Platform{platform}}
		client.ListPlatformsReturns(platforms, nil)
		_ = json.Unmarshal([]byte(args[1]), platform)
		client.UpdatePlatformReturns(&platform, nil)
		ubCmd := command.Prepare(cmd.SmPrepare)
		ubCmd.SetArgs(args)
		return ubCmd.Execute()
	}

	invalidUpdatePlatformExecution := func(args ...string) error {
		ubCmd := command.Prepare(cmd.SmPrepare)
		ubCmd.SetArgs(args)
		return ubCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided", func() {
			It("platform should be updated", func() {

				err := validUpdatePlatformExecution("platform", `{"type":"newType"}`)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring(platform.TableData().String()))
			})

			It("argument values should be as expected", func() {
				err := validUpdatePlatformExecution("platform", `{"type":"newType"}`)

				id, platform := client.UpdatePlatformArgsForCall(0)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(id).To(Equal("id"))
				Expect(platform).To(Equal(&types.Platform{Type: "newType"}))
			})
		})

		Context("With json format flag", func() {
			It("should be printed in json format", func() {
				err := validUpdatePlatformExecution("platform", `{"type":"newType"}`, "--output", "json")

				jsonByte, _ := json.MarshalIndent(platform, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml format flag", func() {
			It("should be printed in yaml format", func() {
				err := validUpdatePlatformExecution("platform", `{"type":"newType"}`, "--output", "yaml")

				yamlByte, _ := yaml.Marshal(platform)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With missing arguments", func() {
			It("Should return error missing name", func() {
				err := invalidUpdatePlatformExecution([]string{}...)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("[name] is required"))
			})
			It("Should return error missing json", func() {
				err := invalidUpdatePlatformExecution("platform")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nothing to update. Platform JSON is not provided"))
			})
		})
	})

	Context("When non existing platform updated", func() {
		It("should return error", func() {
			client.ListPlatformsReturns(&types.Platforms{}, nil)

			err := invalidUpdatePlatformExecution("platform", `{"type":"newType"}`)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("platform with name platform not found"))
		})
	})

	Context("With error from http client", func() {
		It("Should return error", func() {
			expectedErr := errors.New("http client error")
			platforms := &types.Platforms{Platforms: []types.Platform{platform}}
			client.ListPlatformsReturns(platforms, nil)
			client.UpdatePlatformReturns(nil, expectedErr)

			err := invalidUpdatePlatformExecution("platform", `{"type":"newType"}`)

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expectedErr.Error()))
		})
	})

	Context("With invalid output format", func() {
		It("should return error", func() {
			invFormat := "invalid-format"
			err := invalidUpdatePlatformExecution("platform", `{"type":"newType"}`, "--output", invFormat)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown output: " + invFormat))
		})
	})

	Context("when json is invalid", func() {
		It("should throw error", func() {
			err := invalidUpdatePlatformExecution("platform", "{name: none}")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("platform JSON is invalid"))
		})
	})

})
