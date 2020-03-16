package offering

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Marketplace command test", func() {
	var client *smclientfakes.FakeClient
	var command *MarketplaceCmd
	var buffer *bytes.Buffer

	plan1 := types.ServicePlan{
		Name:        "plan1",
		Description: "desc",
	}

	plan2 := types.ServicePlan{
		Name:        "plan2",
		Description: "desc",
	}

	noPlanOffering := types.ServiceOffering{
		Name:        "no-plan-offering",
		Plans:       []types.ServicePlan{},
		Description: "desc",
		BrokerName:  "broker",
		BrokerID:    "id",
	}

	offering1 := types.ServiceOffering{
		Name:        "offering1",
		Plans:       []types.ServicePlan{plan1},
		Description: "desc",
		BrokerName:  "broker1",
		BrokerID:    "id1",
	}

	offering2 := types.ServiceOffering{
		Name:        "offering2",
		Plans:       []types.ServicePlan{plan1, plan2},
		Description: "desc",
		BrokerName:  "broker2",
		BrokerID:    "id2",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewMarketplaceCmd(context)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no offerings provided", func() {
		It("should list empty offerings list", func() {
			client.MarketplaceReturns(&types.Marketplace{ServiceOfferings: []types.ServiceOffering{}}, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("There are no service offerings"))
		})
	})

	Context("when offerings are provided", func() {
		It("should list 1 offering", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{offering1}}
			client.MarketplaceReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})

		It("should list more offerings", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{offering1, offering2}}
			client.MarketplaceReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.Message()))
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})
	})

	Context("when generic parameter is used", func() {
		It("should pass it to SM", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{noPlanOffering}}
			client.MarketplaceReturns(result, nil)
			param := "parameterKey=parameterValue"
			err := executeWithArgs([]string{"--param", param})
			Expect(err).ShouldNot(HaveOccurred())

			args := client.MarketplaceArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when service flag is used", func() {
		It("should list empty plans list when no plans provided", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{noPlanOffering}}
			client.MarketplaceReturns(result, nil)
			err := executeWithArgs([]string{"-s", "no-plan-offering"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("There are no service plans for this service offering."))
		})

		It("should list 1 plan when a single plan is provided", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{offering1}}
			client.MarketplaceReturns(result, nil)
			err := executeWithArgs([]string{"-s", "offering1"})

			expected := &types.ServicePlansForOffering{ServicePlans: result.ServiceOfferings[0].Plans}

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(expected.TableData().String()))
		})

		It("should list multiple plans when multiple plans are provided", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{offering2}}
			client.MarketplaceReturns(result, nil)
			err := executeWithArgs([]string{"-s", "offering2"})

			expected := &types.ServicePlansForOffering{ServicePlans: result.ServiceOfferings[0].Plans}

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(expected.Message()))
			Expect(buffer.String()).To(ContainSubstring(expected.TableData().String()))
		})
	})

	Context("when format flag is used", func() {
		It("should print offerings in json", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{offering1}}
			client.MarketplaceReturns(result, nil)

			err := executeWithArgs([]string{"-o", "json"})

			jsonByte, _ := json.MarshalIndent(result, "", "  ")
			jsonOutputExpected := string(jsonByte) + "\n"
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
		})

		It("should print offerings in yaml", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{offering1}}
			client.MarketplaceReturns(result, nil)

			err := executeWithArgs([]string{"-o", "yaml"})

			yamlByte, _ := yaml.Marshal(result)
			yamlOutputExpected := string(yamlByte) + "\n"
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(yamlOutputExpected))
		})

		It("should print plans in json", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{offering1}}
			client.MarketplaceReturns(result, nil)

			err := executeWithArgs([]string{"-s", "offering1", "-o", "json"})

			jsonByte, _ := json.MarshalIndent(&types.ServicePlansForOffering{ServicePlans: result.ServiceOfferings[0].Plans}, "", "  ")
			jsonOutputExpected := string(jsonByte) + "\n"
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
		})

		It("should print plans in yaml", func() {
			result := &types.Marketplace{ServiceOfferings: []types.ServiceOffering{offering1}}
			client.MarketplaceReturns(result, nil)

			err := executeWithArgs([]string{"-s", "offering1", "-o", "yaml"})

			yamlByte, _ := yaml.Marshal(&types.ServicePlansForOffering{ServicePlans: result.ServiceOfferings[0].Plans})
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
			client.MarketplaceReturns(nil, expectedErr)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expectedErr))
		})
	})

})
