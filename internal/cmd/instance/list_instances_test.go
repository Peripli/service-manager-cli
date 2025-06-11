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
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("List instances command test", func() {

	var client *smclientfakes.FakeClient
	var command *ListInstancesCmd
	var buffer *bytes.Buffer

	instance1 := types.ServiceInstance{
		ID:            "id1",
		Name:          "instance1",
		ServicePlanID: "service_plan_id1",
		PlatformID:    "platform_id1",
	}

	instance2 := types.ServiceInstance{
		ID:            "id2",
		Name:          "instance2",
		ServicePlanID: "service_plan_id2",
		PlatformID:    "platform_id2",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewListInstancesCmd(context)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no instances are registered", func() {
		It("should list empty instances", func() {
			client.ListInstancesReturns(&types.ServiceInstances{ServiceInstances: []types.ServiceInstance{}}, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("There are no service instances."))
		})
	})

	Context("when instances are registered", func() {
		It("should list 1 instance", func() {
			result := &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance1}}
			client.ListInstancesReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})

		It("should list more instances", func() {
			result := &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance1, instance2}}
			client.ListInstancesReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})
	})

	Context("when generic parameter is used", func() {
		It("should pass it to SM", func() {
			result := &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance1}}
			client.ListInstancesReturns(result, nil)
			param := "parameterKey=parameterValue"
			err := executeWithArgs([]string{"--param", param})
			Expect(err).ShouldNot(HaveOccurred())

			args := client.ListInstancesArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when field query flag is used", func() {
		It("should pass it to SM", func() {
			result := &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance1}}
			client.ListInstancesReturns(result, nil)
			param := "name eq 'instance1'"
			err := executeWithArgs([]string{"--field-query", param})

			args := client.ListInstancesArgsForCall(0)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(args.FieldQuery).To(ConsistOf(param))
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when label query flag is used", func() {
		It("should pass it to SM", func() {
			result := &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance1}}
			client.ListInstancesReturns(result, nil)
			param := "test eq false"
			err := executeWithArgs([]string{"--label-query", param})

			args := client.ListInstancesArgsForCall(0)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(ConsistOf(param))
		})
	})

	Context("when format flag is used", func() {
		It("should print in json", func() {
			result := &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance1}}
			client.ListInstancesReturns(result, nil)

			err := executeWithArgs([]string{"-o", "json"})

			jsonByte, _ := json.MarshalIndent(result, "", "  ")
			jsonOutputExpected := string(jsonByte) + "\n"
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
		})

		It("should print in yaml", func() {
			result := &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance1}}
			client.ListInstancesReturns(result, nil)

			err := executeWithArgs([]string{"-o", "yaml"})

			yamlByte, _ := yaml.Marshal(result)
			yamlOutputExpected := string(yamlByte) + "\n"
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(yamlOutputExpected))
		})
	})

	Context("when invalid flag is used", func() {
		It("should handle cobra error", func() {
			err := executeWithArgs([]string{"--ooutput", "json"})
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("unknown flag: --ooutput"))
		})

		It("should handle wrong value", func() {
			err := executeWithArgs([]string{"--output", "invalid"})
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("unknown output: invalid"))
		})
	})

	Context("when error is returned by Service manager", func() {
		It("should handle error", func() {
			expectedErr := errors.New("Http Client Error")
			client.ListInstancesReturns(nil, expectedErr)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expectedErr))
		})
	})
})
