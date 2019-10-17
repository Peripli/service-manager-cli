package visibility

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

func TestListVisibilitiesCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("List visibilities command test", func() {

	var client *smclientfakes.FakeClient
	var command *ListVisibilitiesCmd
	var buffer *bytes.Buffer

	visibility := types.Visibility{
		ID:            "visibilityID",
		PlatformID:    "platformID",
		ServicePlanID: "planID",
	}

	visibility2 := types.Visibility{
		ID:            "visibilityID2",
		PlatformID:    "platformID2",
		ServicePlanID: "planID2",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewListVisibilitiesCmd(context)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no visibilities are registered", func() {
		It("should list empty visibilities", func() {
			client.ListVisibilitiesReturns(&types.Visibilities{Visibilities: []types.Visibility{}}, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("No visibilities registered"))
		})
	})

	Context("when visibilities are registered", func() {
		It("should list 1 visibility", func() {
			result := &types.Visibilities{Visibilities: []types.Visibility{visibility}}
			client.ListVisibilitiesReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})

		It("should list more visibilities", func() {
			result := &types.Visibilities{Visibilities: []types.Visibility{visibility, visibility2}}
			client.ListVisibilitiesReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.Message()))
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})
	})

	Context("when generic parameter is used", func() {
		It("should pass it to SM", func() {
			result := &types.Visibilities{Visibilities: []types.Visibility{visibility}}
			client.ListVisibilitiesReturns(result, nil)
			param := "parameterKey=parameterValue"
			err := executeWithArgs([]string{"--param", param})
			Expect(err).ShouldNot(HaveOccurred())

			args := client.ListVisibilitiesArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when field query flag is used", func() {
		It("should pass it to SM", func() {
			result := &types.Visibilities{Visibilities: []types.Visibility{visibility}}
			client.ListVisibilitiesReturns(result, nil)
			err := executeWithArgs([]string{"-f", "planId eq 'plan1'"})

			queryArg := client.ListVisibilitiesArgsForCall(0)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(queryArg.FieldQuery[0]).To(Equal("planId eq 'plan1'"))
			Expect(queryArg.LabelQuery).To(BeEmpty())
		})
	})

	Context("when label query flag is used", func() {
		It("should pass it to SM", func() {
			result := &types.Visibilities{Visibilities: []types.Visibility{visibility}}
			client.ListVisibilitiesReturns(result, nil)
			err := executeWithArgs([]string{"-l", "test eq false"})

			queryArg := client.ListVisibilitiesArgsForCall(0)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(queryArg.FieldQuery).To(BeEmpty())
			Expect(queryArg.LabelQuery[0]).To(Equal("test eq false"))
		})
	})

	Context("when format flag is used", func() {
		It("should print in json", func() {
			result := &types.Visibilities{Visibilities: []types.Visibility{visibility}}
			client.ListVisibilitiesReturns(result, nil)

			err := executeWithArgs([]string{"-o", "json"})

			jsonByte, _ := json.MarshalIndent(result, "", "  ")
			jsonOutputExpected := string(jsonByte) + "\n"
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
		})

		It("should print in yaml", func() {
			result := &types.Visibilities{Visibilities: []types.Visibility{visibility}}
			client.ListVisibilitiesReturns(result, nil)

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
			client.ListVisibilitiesReturns(nil, expectedErr)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expectedErr))
		})
	})

})
