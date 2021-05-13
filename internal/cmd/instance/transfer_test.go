/*
 * Copyright 2018 The Service Manager Authors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package instance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Transfer Command test", func() {
	var client *smclientfakes.FakeClient
	var command *TransferCmd
	var buffer *bytes.Buffer
	var promptBuffer *bytes.Buffer

	var instance *types.ServiceInstance

	validAsyncTransferExecution := func(location string, args ...string) *cobra.Command {
		instance = &types.ServiceInstance{
			Name: args[0],
		}
		operation := &types.Operation{
			State: "in progress",
		}
		client.StatusReturns(operation, nil)
		client.ListInstancesReturns(&types.ServiceInstances{ServiceInstances: []types.ServiceInstance{*instance}}, nil)
		client.UpdateInstanceReturns(instance, location, nil)

		piCmd := command.Prepare(cmd.SmPrepare)
		piCmd.SetArgs(args)
		Expect(piCmd.Execute()).ToNot(HaveOccurred())

		return piCmd
	}

	validSyncTransferExecExpect := func(args ...string) *cobra.Command {
		return validAsyncTransferExecution("", append(args, "--mode", "sync")...)
	}

	invalidTransferCommandExecution := func(args ...string) error {
		trCmd := command.Prepare(cmd.SmPrepare)
		trCmd.SetArgs(args)
		return trCmd.Execute()
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		promptBuffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewTransferCmd(context, promptBuffer)
	})

	Describe("Valid request", func() {
		BeforeEach(func() {
			promptBuffer.WriteString("y")
		})

		Context("With necessary arguments provided", func() {
			It("should be transfered successfully", func() {
				validSyncTransferExecExpect("instance-name", "--from", "platform_id", "--to", "service-manager")

				tableOutputExpected := instance.TableData().String()

				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})

			It("should print location when transfered asynchronously", func() {
				validAsyncTransferExecution("location", "instance-name", "--from", "platform_id", "--to", "service-manager")

				Expect(buffer.String()).To(ContainSubstring(`smctl status location`))
			})

			It("Argument values should be as expected", func() {
				validSyncTransferExecExpect("instance-name", "--from", "from_platform", "--to", "to_platform")

				Expect(command.instanceName).To(Equal("instance-name"))
				Expect(command.fromPlatformID).To(Equal("from_platform"))
				Expect(command.toPlatformID).To(Equal("to_platform"))
			})
		})

		Context("when 2 instances are present with same name", func() {
			BeforeEach(func() {
				client.ListInstancesReturnsOnCall(0, &types.ServiceInstances{
					ServiceInstances: []types.ServiceInstance{
						types.ServiceInstance{
							Name: "instance-name",
						},
						types.ServiceInstance{
							Name: "instance-name",
						},
					},
				}, nil)

			})
			Context("when no instance id is provided", func() {
				It("should require flag for instance id", func() {
					err := invalidTransferCommandExecution("instance-name", "--from", "from_platform", "--to", "to_platform")
					Expect(err.Error()).To(Equal(fmt.Sprintf(cmd.FOUND_TOO_MANY_INSTANCES,"instance-name","transfer")))
				})
			})

			Context("when instance id is provided", func() {
				It("should transfer the specified instance id", func() {
					validSyncTransferExecExpect("instance-name", "--from", "from_platform", "--to", "to_platform", "--id", "12345")
					Expect(buffer.String()).To(ContainSubstring(instance.TableData().String()))
				})
			})
		})

		Context("when no instanes are present with certain name", func() {
			BeforeEach(func() {
				client.ListInstancesReturnsOnCall(0, &types.ServiceInstances{
					ServiceInstances: []types.ServiceInstance{},
				}, nil)
			})

			It("should fail to transfer", func() {
				err := invalidTransferCommandExecution("no-instance", "--from", "from_platform", "--to", "to_platform")
				message:=fmt.Sprintf(cmd.NO_INSTANCES_FOUND,"no-instance")
				Expect(err.Error()).To(Equal(message))
			})
		})

		Context("With json output flag", func() {
			It("should be printed in json output format", func() {
				validSyncTransferExecExpect("instance-name", "--from", "from_platform", "--to", "to_platform", "--output", "json")

				jsonByte, _ := json.MarshalIndent(instance, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
			})
		})

		Context("With yaml output flag", func() {
			It("should be printed in yaml output format", func() {
				validSyncTransferExecExpect("instance-name", "--from", "from_platform", "--to", "to_platform", "--output", "yaml")

				yamlByte, _ := yaml.Marshal(instance)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(buffer.String()).To(ContainSubstring(yamlOutputExpected))
			})
		})
	})

	Describe("Invalid requests", func() {
		BeforeEach(func() {
			promptBuffer.WriteString("y")
		})

		When("list instances fails", func() {
			BeforeEach(func() {
				client.ListInstancesReturns(nil, errors.New("errored"))
			})

			It("should return error", func() {
				err := invalidTransferCommandExecution("instance-name", "--from", "from_platform", "--to", "to_platform")
				Expect(err.Error()).To(Equal("errored"))
			})
		})

		When("update instance fails", func() {
			BeforeEach(func() {
				client.ListInstancesReturnsOnCall(0, &types.ServiceInstances{
					ServiceInstances: []types.ServiceInstance{
						types.ServiceInstance{
							Name: "instance-name",
						},
						types.ServiceInstance{
							Name: "instance-name",
						},
					},
				}, nil)

				client.UpdateInstanceReturns(nil, "", errors.New("errored"))
			})

			It("should return error", func() {
				err := invalidTransferCommandExecution("instance-name", "--from", "from_platform", "--to", "to_platform", "--id", "1234")
				Expect(err.Error()).To(Equal("errored"))
			})
		})

		When("when transfer is declined", func() {
			BeforeEach(func() {
				promptBuffer.Reset()
				promptBuffer.WriteString("n")
			})

			It("should print appropriate message", func() {
				validSyncTransferExecExpect("instance-name", "--from", "from_platform", "--to", "to_platform", "--id", "1234")

				Expect(buffer.String()).To(ContainSubstring("Transfer declined"))
			})
		})
	})

})
