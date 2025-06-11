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
	"fmt"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/Peripli/service-manager/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Update sharing instance command test", func() {
	var client *smclientfakes.FakeClient
	var command *UpdateSharingCmd
	var buffer *bytes.Buffer
	var instance *types.ServiceInstance
	var instances *types.ServiceInstances
	validUpdateSharingExecution := func(args ...string) *cobra.Command {
		instance = &types.ServiceInstance{
			Name: args[0],
		}
		for i, arg := range args {
			if arg == "--id" {
				instance.ID = args[i+1]
			}
		}
		operation := &types.Operation{
			State: "in progress",
		}

		client.StatusReturns(operation, nil)
		client.ListInstancesReturns(&types.ServiceInstances{ServiceInstances: []types.ServiceInstance{*instance}}, nil)
		client.UpdateInstanceReturns(instance, "", nil)
		piCmd := command.Prepare(cmd.SmPrepare)
		piCmd.SetArgs(args)
		Expect(piCmd.Execute()).ToNot(HaveOccurred())
		return piCmd
	}

	invalidUpdateSharingCommandExecution := func(args ...string) error {
		shCmd := command.Prepare(cmd.SmPrepare)
		shCmd.SetArgs(args)
		return shCmd.Execute()
	}
	type testCase struct {
		share       bool
		commandName string
	}
	tests := []testCase{
		testCase{true, "share"},
		testCase{false, "unshare"},
	}
	for _, test := range tests {
		Describe(fmt.Sprintf("%s instance", test.commandName), func() {
			BeforeEach(func() {
				buffer = &bytes.Buffer{}
				client = &smclientfakes.FakeClient{}
				context := &cmd.Context{Output: buffer, Client: client}
				command = NewUpdateSharingCmd(context, test.share)
			})

			Context("valid sync", func() {
				Context("with name", func() {
					It("should print object", func() {
						validUpdateSharingExecution("myinstancename")
						tableOutputExpected := instance.TableData().String()
						Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
					})

				})
				Context("with id", func() {
					It("should print object", func() {
						validUpdateSharingExecution("myinstancename", "--id", "serviceinstanceid")
						tableOutputExpected := instance.TableData().String()
						Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
					})
				})

				Context("with json output flag", func() {
					It("should be printed in json output format", func() {
						validUpdateSharingExecution("instance-name", "--output", "json")
						jsonByte, _ := json.MarshalIndent(instance, "", "  ")
						jsonOutputExpected := string(jsonByte) + "\n"
						Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
					})
				})

				Context("with yaml output flag", func() {
					It("should be printed in yaml output format", func() {
						validUpdateSharingExecution("instance-name", "--output", "yaml")
						yamlByte, _ := yaml.Marshal(instance)
						yamlOutputExpected := string(yamlByte) + "\n"
						Expect(buffer.String()).To(ContainSubstring(yamlOutputExpected))
					})
				})
			})

			Context("invalid execution", func() {
				JustBeforeEach(func() {
					client.ListInstancesReturns(instances, nil)

				})

				When("invalid flag", func() {
					It("should return an error", func() {
						err := invalidUpdateSharingCommandExecution("instance name", "--fl", "fff")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("unknown flag: --fl"))
					})
				})
				When("async is used", func() {
					It("should return an error", func() {
						err := invalidUpdateSharingCommandExecution("instance name", "--mode", "async")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("unknown flag: --mode"))
					})
				})

				When("service instance not found", func() {
					BeforeEach(func() {
						instances = &types.ServiceInstances{}
					})

					It("should return an error", func() {
						err := invalidUpdateSharingCommandExecution("instance-name")
						Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(cmd.NO_INSTANCES_FOUND, "instance-name")))
					})

				})
				Context("when more than once instance found", func() {
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
					It("should return an error", func() {
						err := invalidUpdateSharingCommandExecution("instance-name")
						Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(cmd.FOUND_TOO_MANY_INSTANCES, "instance-name", test.commandName)))
					})

				})
				Context("backend error", func() {
					BeforeEach(func() {
						client.ListInstancesReturnsOnCall(0, &types.ServiceInstances{
							ServiceInstances: []types.ServiceInstance{
								types.ServiceInstance{
									Name: "instance-name",
								},
							},
						}, nil)
						body := ioutil.NopCloser(bytes.NewReader([]byte("HTTP response error")))
						expectedError := util.HandleResponseError(&http.Response{Body: body})
						client.UpdateInstanceReturns(nil, "", expectedError)
					})
					It("should return error's description", func() {
						err := invalidUpdateSharingCommandExecution("instance-name")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("HTTP response error"))
						Expect(buffer.String()).To(ContainSubstring(fmt.Sprintf("Couldn't %s the service instance. ", test.commandName)))

					})

				})

				Context("with invalid output format", func() {
					It("should return error", func() {
						invFormat := "invalid-format"
						err := invalidUpdateSharingCommandExecution("instance name", "--output", invFormat)
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(Equal("unknown output: " + invFormat))
					})
				})

			})

		})

	}
})
