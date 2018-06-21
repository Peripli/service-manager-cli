package broker

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"bytes"
	"errors"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	resperrors "github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/spf13/cobra"
)

func TestRegisterBrokerCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Register Broker Command test", func() {

	var client *smclientfakes.FakeClient
	var command *RegisterBrokerCmd
	var buffer *bytes.Buffer
	var broker *types.Broker

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewRegisterBrokerCmd(context)
	})

	validRegisterBrokerExecution := func(args []string) *cobra.Command {
		broker = &types.Broker{
			Name: args[0],
			URL:  args[1],
			ID:   "1234",
		}
		client.RegisterBrokerReturns(broker, nil)

		rbcCmd := command.Prepare(cmd.SmPrepare)
		rbcCmd.SetArgs(args)
		rbcCmd.Execute()

		return rbcCmd
	}

	invalidRegisterBrokerCommandExecution := func(args []string) error {
		rpcCmd := command.Prepare(cmd.SmPrepare)
		rpcCmd.SetArgs(args)
		return rpcCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided and basic flag", func() {
			It("should be registered", func() {
				validRegisterBrokerExecution([]string{"broker-name", "http://broker.com", "--basic", "user:password"})

				tableOutputExpected := broker.TableData().String()

				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})

			It("Argument values should be as expected", func() {
				validRegisterBrokerExecution([]string{"broker-name", "http://broker.com", "--basic", "user:password"})

				Expect(command.broker.Name).To(Equal("broker-name"))
				Expect(command.broker.URL).To(Equal("http://broker.com"))
				Expect(command.broker.Credentials.Basic.User).To(Equal("user"))
				Expect(command.broker.Credentials.Basic.Password).To(Equal("password"))
			})
		})

		Context("With description provided", func() {
			It("should save description value as expected", func() {
				validRegisterBrokerExecution([]string{"validName", "validType", "validDescription", "--basic", "user:password"})

				Expect(command.broker.Description).To(Equal("validDescription"))
			})
		})

		Context("With json format flag", func() {
			It("should be printed in json format", func() {
				validRegisterBrokerExecution([]string{"validName", "validUrl", "--basic", "user:password", "--format", "json"})

				jsonByte, _ := json.MarshalIndent(broker, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml format flag", func() {
			It("should be printed in yaml format", func() {
				validRegisterBrokerExecution([]string{"validName", "validUrl", "--basic", "user:password", "--format", "yaml"})

				yamlByte, _ := yaml.Marshal(broker)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With not enough arguments provided", func() {
			It("should return error", func() {
				err := invalidRegisterBrokerCommandExecution([]string{"validName", "--basic", "user:password"})

				Expect(err.Error()).To(ContainSubstring("Name and URL are required"))
			})
		})

		Context("With invalid basic flag provided", func() {
			It("should return error", func() {
				err := invalidRegisterBrokerCommandExecution([]string{"validName", "validType", "--basic", "invalidBasicFlag"})

				Expect(err.Error()).To(ContainSubstring("basic string is invalid"))
			})
		})

		Context("With error from http client", func() {
			It("should return error", func() {
				client.RegisterBrokerReturns(nil, errors.New("Http Client Error"))

				err := invalidRegisterBrokerCommandExecution([]string{"validName", "validType", "--basic", "user:password"})

				Expect(err).To(MatchError("Http Client Error"))
			})
		})

		Context("With http response error from http client", func() {
			It("should return error's description", func() {
				client.RegisterBrokerReturns(nil, resperrors.ResponseError{Description: "HTTP response error"})

				err := invalidRegisterBrokerCommandExecution([]string{"validName", "validType", "--basic", "user:password"})

				Expect(err).To(MatchError("HTTP response error"))
			})
		})

		Context("With invalid output format", func() {
			It("should return error", func() {
				invFormat := "invalid-format"
				err := invalidRegisterBrokerCommandExecution([]string{"validName", "validUrl", "--basic", "user:password", "--format", invFormat})

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("unknown format: " + invFormat))
			})
		})
	})
})
