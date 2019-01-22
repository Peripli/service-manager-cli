package platform

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"

	"bytes"

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
	platform := types.Platform{
		Name: "platform",
		ID:   "id1",
		Type: "type1",
		Credentials: &types.Credentials{
			Basic: types.Basic{
				User:     "admin",
				Password: "admin",
			},
		},
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewUpdatePlatformCmd(context)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing platform is being updated", func() {
		It("should print updated platform", func() {
			updatedPlatform := platform
			updatedPlatform.Name = "updated-name"
			client.UpdatePlatformReturns(&updatedPlatform, nil)
			client.ListPlatformsReturns(&types.Platforms{Platforms: []types.Platform{platform}}, nil)
			err := executeWithArgs([]string{platform.Name, `{"name": "` + updatedPlatform.Name + `"}`})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(updatedPlatform.TableData().String()))
		})
	})

	Context("when non-existing platform is being updated", func() {
		It("should throw error", func() {
			client.ListPlatformsReturns(&types.Platforms{Platforms: []types.Platform{platform}}, nil)
			err := executeWithArgs([]string{"non-existing", "{}"})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("platform with name non-existing not found"))
		})
	})

	Context("when list platforms returns an error", func() {
		It("should be handled", func() {
			client.ListPlatformsReturns(nil, errors.New("error retrieving platforms"))
			err := executeWithArgs([]string{"non-existing", "{}"})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("error retrieving platforms"))
		})
	})

	Context("when name is not provided", func() {
		It("should throw error", func() {
			err := executeWithArgs([]string{})
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("[name] is required"))
		})
	})

	Context("when json is not provided", func() {
		It("should throw error", func() {
			err := executeWithArgs([]string{"platform"})
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("nothing to update. Platform JSON is not provided"))
		})
	})

	Context("when json is invalid", func() {
		It("should throw error", func() {
			err := executeWithArgs([]string{"platform", "{name: none}"})
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("platform JSON is invalid"))
		})
	})

	Context("when format flag is used", func() {
		It("should print in json", func() {
			client.UpdatePlatformReturns(&platform, nil)
			client.ListPlatformsReturns(&types.Platforms{Platforms: []types.Platform{platform}}, nil)
			executeWithArgs([]string{platform.Name, `{"name": "platform"}`, "-o", "json"})

			jsonByte, _ := json.MarshalIndent(platform, "", "  ")
			jsonOutputExpected := string(jsonByte) + "\n"
			Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
		})

		It("should print in yaml", func() {
			client.UpdatePlatformReturns(&platform, nil)
			client.ListPlatformsReturns(&types.Platforms{Platforms: []types.Platform{platform}}, nil)
			executeWithArgs([]string{platform.Name, `{"name": "platform"}`, "-o", "yaml"})

			yamlByte, _ := yaml.Marshal(platform)
			yamlOutputExpected := string(yamlByte) + "\n"
			Expect(buffer.String()).To(ContainSubstring(yamlOutputExpected))
		})
	})
})
