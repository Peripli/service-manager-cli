package broker

import (
	"encoding/json"
	"github.com/Peripli/service-manager/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"

	"bytes"
	"errors"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/spf13/cobra"
)

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

	validAsyncRegisterBrokerExecution := func(args []string, location string) *cobra.Command {
		broker = &types.Broker{
			Name: args[0],
			URL:  args[1],
			ID:   "1234",
		}
		client.RegisterBrokerReturns(broker, location, nil)

		rbcCmd := command.Prepare(cmd.SmPrepare)
		rbcCmd.SetArgs(args)
		Expect(rbcCmd.Execute()).ToNot(HaveOccurred())

		return rbcCmd
	}

	validSyncRegisterBrokerExecution := func(args []string) *cobra.Command {
		return validAsyncRegisterBrokerExecution(args, "")
	}

	invalidRegisterBrokerCommandExecution := func(args []string) error {
		rpcCmd := command.Prepare(cmd.SmPrepare)
		rpcCmd.SetArgs(args)
		return rpcCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided and basic flag", func() {
			It("should be registered synchronously", func() {
				validSyncRegisterBrokerExecution([]string{"broker-name", "http://broker.com", "--basic", "user:password"})

				tableOutputExpected := broker.TableData().String()

				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})

			It("should print location when registered asynchronously", func() {
				validAsyncRegisterBrokerExecution([]string{"broker-name", "http://broker.com", "--basic", "user:password", "--mode", "async"}, "location")

				Expect(buffer.String()).To(ContainSubstring(`smctl status location`))
			})

			It("Argument values should be as expected", func() {
				validSyncRegisterBrokerExecution([]string{"broker-name", "http://broker.com", "--basic", "user:password"})

				Expect(command.broker.Name).To(Equal("broker-name"))
				Expect(command.broker.URL).To(Equal("http://broker.com"))
				Expect(command.broker.Credentials.Basic.User).To(Equal("user"))
				Expect(command.broker.Credentials.Basic.Password).To(Equal("password"))
			})
		})

		Context("With description provided", func() {
			It("should save description value as expected", func() {
				validSyncRegisterBrokerExecution([]string{"validName", "validType", "validDescription", "--basic", "user:password"})

				Expect(command.broker.Description).To(Equal("validDescription"))
			})
		})

		Context("With json output flag", func() {
			It("should be printed in json output format", func() {
				validSyncRegisterBrokerExecution([]string{"validName", "validUrl", "--basic", "user:password", "--output", "json"})

				jsonByte, _ := json.MarshalIndent(broker, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml output flag", func() {
			It("should be printed in yaml output format", func() {
				validSyncRegisterBrokerExecution([]string{"validName", "validUrl", "--basic", "user:password", "--output", "yaml"})

				yamlByte, _ := yaml.Marshal(broker)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})

		Context("With generic param flag", func() {
			It("should pass it to SM", func() {
				validSyncRegisterBrokerExecution([]string{"validName", "validType", "validDescription", "--basic", "user:password", "--param", "paramKey=paramValue"})

				_, args := client.RegisterBrokerArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("paramKey=paramValue", "async=false"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})

		Context("With async flag", func() {
			It("should pass it to SM", func() {
				validSyncRegisterBrokerExecution([]string{"validName", "validType", "validDescription", "--basic", "user:password", "--mode", "async"})

				_, args := client.RegisterBrokerArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("async=true"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With not enough arguments provided", func() {
			It("should return error", func() {
				err := invalidRegisterBrokerCommandExecution([]string{"validName", "--basic", "user:password"})

				Expect(err.Error()).To(ContainSubstring("name and URL are required"))
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
				client.RegisterBrokerReturns(nil, "", errors.New("Http Client Error"))

				err := invalidRegisterBrokerCommandExecution([]string{"validName", "validType", "--basic", "user:password"})

				Expect(err).To(MatchError("Http Client Error"))
			})
		})

		Context("With http response error from http client", func() {
			It("should return error's description", func() {
				body := ioutil.NopCloser(bytes.NewReader([]byte("HTTP response error")))
				expectedError := util.HandleResponseError(&http.Response{Body: body})
				client.RegisterBrokerReturns(nil, "", expectedError)

				err := invalidRegisterBrokerCommandExecution([]string{"validName", "validType", "--basic", "user:password"})

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("HTTP response error"))
			})
		})

		Context("With invalid output format", func() {
			It("should return error", func() {
				invFormat := "invalid-format"
				err := invalidRegisterBrokerCommandExecution([]string{"validName", "validUrl", "--basic", "user:password", "--output", invFormat})

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("unknown output: " + invFormat))
			})
		})
	})
})
