package platform

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

func TestListPlatformsCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("List platforms command test", func() {

	var client *smclientfakes.FakeClient
	var command *ListPlatformsCmd
	var buffer *bytes.Buffer
	platform := types.Platform{
		Name: "platform1",
		Type: "type1",
		ID:   "id1",
	}
	platform2 := types.Platform{
		Name: "platform2",
		Type: "type2",
		ID:   "id2",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewListPlatformsCmd(context)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no platforms are registered", func() {
		It("should list empty platforms", func() {
			client.ListPlatformsReturns(&types.Platforms{Platforms: []types.Platform{}}, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("No platforms registered"))
		})
	})

	Context("when platforms are registered", func() {
		It("should list 1 platform", func() {
			result := &types.Platforms{Platforms: []types.Platform{platform}}
			client.ListPlatformsReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})

		It("should list more platforms", func() {
			result := &types.Platforms{Platforms: []types.Platform{platform, platform2}}
			client.ListPlatformsReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.Message()))
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})
	})

	Context("when generic parameter is used", func() {
		It("should pass it to SM", func() {
			result := &types.Platforms{Platforms: []types.Platform{platform}}
			client.ListPlatformsReturns(result, nil)
			param := "parameterKey=parameterValue"
			err := executeWithArgs([]string{"--param", param})
			Expect(err).ShouldNot(HaveOccurred())

			args := client.ListPlatformsArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when field query flag is used", func() {
		It("should pass it to SM", func() {
			result := &types.Platforms{Platforms: []types.Platform{platform}}
			client.ListPlatformsReturns(result, nil)
			param := "name = platform1"
			err := executeWithArgs([]string{"--field-query", param})

			args := client.ListPlatformsArgsForCall(0)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(args.FieldQuery).To(ConsistOf(param))
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when label query flag is used", func() {
		It("should pass it to SM", func() {
			result := &types.Platforms{Platforms: []types.Platform{platform}}
			client.ListPlatformsReturns(result, nil)
			param := "test = false"
			err := executeWithArgs([]string{"--label-query", param})

			args := client.ListPlatformsArgsForCall(0)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(args.LabelQuery).To(ConsistOf(param))
			Expect(args.FieldQuery).To(BeEmpty())
		})
	})

	Context("when format flag is used", func() {
		It("should print in json", func() {
			result := &types.Platforms{Platforms: []types.Platform{platform}}
			client.ListPlatformsReturns(result, nil)

			err := executeWithArgs([]string{"-o", "json"})

			jsonByte, _ := json.MarshalIndent(result, "", "  ")
			jsonOutputExpected := string(jsonByte) + "\n"
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
		})

		It("should print in yaml", func() {
			result := &types.Platforms{Platforms: []types.Platform{platform}}
			client.ListPlatformsReturns(result, nil)

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
			client.ListPlatformsReturns(nil, expectedErr)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expectedErr))
		})
	})

})
