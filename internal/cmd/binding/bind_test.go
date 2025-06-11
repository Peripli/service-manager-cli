package binding

import (
	"encoding/json"
	"github.com/Peripli/service-manager/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"

	"bytes"
	"errors"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/spf13/cobra"
)

var _ = Describe("Bind Command test", func() {
	var client *smclientfakes.FakeClient
	var command *BindCmd
	var buffer *bytes.Buffer

	var binding *types.ServiceBinding
	var instance *types.ServiceInstance

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewBindCmd(context)
	})

	validAsyncBindExecution := func(location string, args ...string) *cobra.Command {
		instance = &types.ServiceInstance{
			ID:   "instance-id",
			Name: args[0],
		}
		binding = &types.ServiceBinding{
			ID:   "instance-id",
			Name: args[1],
		}
		operation := &types.Operation{
			State: "in progress",
		}
		client.StatusReturns(operation, nil)
		client.ListInstancesReturns(&types.ServiceInstances{ServiceInstances: []types.ServiceInstance{*instance}}, nil)
		client.BindReturns(binding, location, nil)

		bindCmd := command.Prepare(cmd.SmPrepare)
		bindCmd.SetArgs(args)
		Expect(bindCmd.Execute()).ToNot(HaveOccurred())

		return bindCmd
	}

	validSyncBindExecution := func(args ...string) *cobra.Command {
		return validAsyncBindExecution("", append(args, "--mode", "sync")...)
	}

	invalidBindCommandExecution := func(args ...string) error {
		rpcCmd := command.Prepare(cmd.SmPrepare)
		rpcCmd.SetArgs(args)
		return rpcCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided", func() {
			It("should be registered synchronously", func() {
				validSyncBindExecution("instance-name", "binding-name")

				tableOutputExpected := binding.TableData().String()

				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})

			It("should print location when registered asynchronously", func() {
				validAsyncBindExecution("location", "instance-name", "binding-name")

				Expect(buffer.String()).To(ContainSubstring(`smctl status location`))
			})

			It("Argument values should be as expected", func() {
				validSyncBindExecution("instance-name", "binding-name")

				Expect(command.binding.Name).To(Equal("binding-name"))
				Expect(command.binding.ServiceInstanceID).To(Equal(instance.ID))
			})
		})

		Context("With json output flag", func() {
			It("should be printed in json output format", func() {
				validSyncBindExecution("instance-name", "binding-name", "--output", "json")

				jsonByte, _ := json.MarshalIndent(binding, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml output flag", func() {
			It("should be printed in yaml output format", func() {
				validSyncBindExecution("instance-name", "binding-name", "--output", "yaml")

				yamlByte, _ := yaml.Marshal(binding)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})

		Context("With generic param flag", func() {
			It("should pass it to SM", func() {
				validSyncBindExecution("instance-name", "binding-name", "--param", "paramKey=paramValue")

				_, args := client.BindArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("paramKey=paramValue", "async=false"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})

		Context("With sync flag", func() {
			It("should pass it to SM", func() {
				validSyncBindExecution("instance-name", "binding-name")

				_, args := client.BindArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("async=false"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})
	})

	Describe("Invalid request", func() {
		var instances *types.ServiceInstances

		BeforeEach(func() {
			instances = &types.ServiceInstances{
				ServiceInstances: []types.ServiceInstance{
					{ID: "id", Name: "instance-name"},
				},
			}
		})

		JustBeforeEach(func() {
			client.ListInstancesReturns(instances, nil)
		})

		Context("With not enough arguments provided", func() {
			It("should return error", func() {
				err := invalidBindCommandExecution("instance-name")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("instance and binding names are required"))
			})
		})

		Context("When instance not found", func() {
			BeforeEach(func() {
				instances = &types.ServiceInstances{}
			})

			It("should return error", func() {
				err := invalidBindCommandExecution("instance-name", "binding-name")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(SatisfyAll(ContainSubstring("service instance with name"), ContainSubstring("not found")))
			})
		})

		Context("when more than one instance with given name found", func() {
			BeforeEach(func() {
				instances.ServiceInstances = append(instances.ServiceInstances, types.ServiceInstance{ID: "456", Name: "instance-name"})
			})
			It("should return message", func() {
				err := invalidBindCommandExecution("instance-name", "binding-name")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(SatisfyAll(ContainSubstring("more than one service instance with name"), ContainSubstring("found. Use --id flag to specify id of the instance to bind")))
			})
		})

		Context("With error from http client", func() {
			It("should return error", func() {
				client.BindReturns(nil, "", errors.New("Http Client Error"))

				err := invalidBindCommandExecution("instance-name", "binding-name")

				Expect(err).To(MatchError("Http Client Error"))
			})
		})

		Context("With http response error from http client", func() {
			It("should return error's description", func() {
				body := ioutil.NopCloser(bytes.NewReader([]byte("HTTP response error")))
				expectedError := util.HandleResponseError(&http.Response{Body: body})
				client.BindReturns(nil, "", expectedError)

				err := invalidBindCommandExecution("instance-name", "binding-name")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("HTTP response error"))
			})
		})

		Context("With invalid output format", func() {
			It("should return error", func() {
				invFormat := "invalid-format"
				err := invalidBindCommandExecution("instance-name", "binding-name", "--output", invFormat)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("unknown output: " + invFormat))
			})
		})
	})
})
