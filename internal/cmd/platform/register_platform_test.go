package platform

import (
	"encoding/json"
	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"errors"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

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

	validRegisterPlatformExecution := func(args ...string) error {
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

		rpcCmd := command.Prepare(cmd.SmPrepare)
		rpcCmd.SetArgs(args)
		return rpcCmd.Execute()
	}

	invalidRegisterPlatformCommandExecution := func(args ...string) error {
		rpcCmd := command.Prepare(cmd.SmPrepare)
		rpcCmd.SetArgs(args)
		return rpcCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided", func() {
			It("Platform should be registered", func() {
				err := validRegisterPlatformExecution("platform", "cf")
				tableOutputExpected := platform.TableData().String()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})

			It("Argument values should be as expected", func() {
				err := validRegisterPlatformExecution("validName", "validType")

				p, _ := client.RegisterPlatformArgsForCall(0)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(p.Name).To(Equal("validName"))
				Expect(p.Type).To(Equal("validType"))
			})
		})

		Context("With id flag provided", func() {
			It("Flag value should be as expected", func() {
				args := []string{"validName", "validType", "--id", "1234"}

				err := validRegisterPlatformExecution(args...)
				platform, _ := client.RegisterPlatformArgsForCall(0)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(platform.ID).To(Equal("1234"))
			})
		})

		Context("With description provided", func() {
			It("Description value should be as expected", func() {
				err := validRegisterPlatformExecution("validName", "validType", "validDescription")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(command.platform.Description).To(Equal("validDescription"))
			})
		})

		Context("With json format flag", func() {
			It("should be printed in json format", func() {
				err := validRegisterPlatformExecution("platform", "cf", "--output", "json")

				jsonByte, _ := json.MarshalIndent(platform, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml format flag", func() {
			It("should be printed in yaml format", func() {
				err := validRegisterPlatformExecution("platform", "cf", "--output", "yaml")

				yamlByte, _ := yaml.Marshal(platform)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})

		Context("With generic param flag", func() {
			It("should pass it to SM", func() {
				err := validRegisterPlatformExecution("platform", "cf", "--param", "paramKey=paramValue")
				Expect(err).ShouldNot(HaveOccurred())

				_, args := client.RegisterPlatformArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("paramKey=paramValue"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With not enough arguments provided", func() {
			It("Should return error", func() {
				err := invalidRegisterPlatformCommandExecution("validName")

				Expect(err.Error()).To(ContainSubstring("requires at least 2 args"))
			})
		})

		Context("With error from http client", func() {
			It("Should return error", func() {
				expectedErr := errors.New("http client error")
				client.RegisterPlatformReturns(nil, expectedErr)

				err := invalidRegisterPlatformCommandExecution("validName", "validType")

				Expect(err).To(MatchError(expectedErr.Error()))
			})
		})

		Context("With invalid output format", func() {
			It("should return error", func() {
				invFormat := "invalid-format"
				err := invalidRegisterPlatformCommandExecution("validName", "validUrl", "--output", invFormat)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("unknown output: " + invFormat))
			})
		})
	})
})
