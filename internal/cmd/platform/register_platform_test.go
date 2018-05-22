package platform

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"errors"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/spf13/cobra"
)

func TestRegisterPlatformCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Register Platform Command test", func() {

	var client *smclientfakes.FakeClient
	var command *RegisterPlatformCmd
	var buffer *bytes.Buffer
	var platform *types.Platform

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewRegisterPlatformCmd(context)
	})

	validRegisterPlatformExecution := func(args []string) *cobra.Command {
		platform = &types.Platform{
			ID:   "1234",
			Name: args[0],
			Type: args[1],
			Credentials: &types.Credentials{
				Basic: types.Basic{
					User:     "admin",
					Password: "admin",
				},
			},
		}
		client.RegisterPlatformReturns(platform, nil)

		rpcCmd := command.Command()
		rpcCmd.SetArgs(args)
		rpcCmd.Execute()

		return rpcCmd
	}

	invalidRegisterPlatformCommandExecution := func(args []string) error {
		rpcCmd := command.Command()
		rpcCmd.SetArgs(args)
		return rpcCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided", func() {
			It("Platform should be registered", func() {
				validRegisterPlatformExecution([]string{"platform", "cf"})
				tableOutputExpected := platform.TableData().String()

				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})

			It("Argument values should be as expected", func() {
				validRegisterPlatformExecution([]string{"validName", "validType"})

				Expect(command.platform.Name).To(Equal("validName"))
				Expect(command.platform.Type).To(Equal("validType"))
			})
		})

		Context("With id flag provided", func() {
			It("Flag value should be as expected", func() {
				args := []string{"validName", "validType"}
				args = append(args, "--id", "1234")

				rpcCmd := validRegisterPlatformExecution(args)
				idFlag := rpcCmd.Flag("id")

				Expect(idFlag.Value.String()).To(Equal("1234"))
				Expect(idFlag.DefValue).To(Equal(""))
			})
		})

		Context("With description provided", func() {
			It("Description value should be as expected", func() {
				validRegisterPlatformExecution([]string{"validName", "validType", "validDescription"})

				Expect(command.platform.Description).To(Equal("validDescription"))
			})
		})

		Context("With json format flag", func() {
			It("should be printed in json format", func() {
				validRegisterPlatformExecution([]string{"platform", "cf", "--format", "json"})

				jsonByte, _ := json.MarshalIndent(platform, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml format flag", func() {
			It("should be printed in yaml format", func() {
				validRegisterPlatformExecution([]string{"platform", "cf", "--format", "yaml"})

				yamlByte, _ := yaml.Marshal(platform)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With not enough arguments provided", func() {
			It("Should return error", func() {
				err := invalidRegisterPlatformCommandExecution([]string{"validName"})

				Expect(err.Error()).To(ContainSubstring("requires at least 2 args"))
			})
		})

		Context("With error from http client", func() {
			It("Should return error", func() {
				client.RegisterPlatformReturns(nil, errors.New("Http Client Error"))

				err := invalidRegisterPlatformCommandExecution([]string{"validName", "validType"})

				Expect(err).To(MatchError("Http Client Error"))
			})
		})

		Context("With invalid output format", func() {
			It("should return error", func() {
				invFormat := "invalid-format"
				err := invalidRegisterPlatformCommandExecution([]string{"validName", "validUrl", "--format", invFormat})

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("unknown format: " + invFormat))
			})
		})
	})
})
