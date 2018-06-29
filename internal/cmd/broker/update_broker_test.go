package broker

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

func TestUpdateBrokerCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Update broker command test", func() {

	var client *smclientfakes.FakeClient
	var command *UpdateBrokerCmd
	var buffer *bytes.Buffer
	broker := types.Broker{
		Name: "broker1",
		ID:   "id1",
		URL:  "http://broker1.com",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewUpdateBrokerCmd(context)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing broker is being updated", func() {
		It("should print updated broker", func() {
			updatedBroker := broker
			updatedBroker.Name = "updated-name"
			client.UpdateBrokerReturns(&updatedBroker, nil)
			client.ListBrokersReturns(&types.Brokers{Brokers: []types.Broker{broker}}, nil)
			err := executeWithArgs([]string{broker.Name, `{"name": "` + updatedBroker.Name + `"}`})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(updatedBroker.TableData().String()))
		})
	})

	Context("when non-existing broker is being updated", func() {
		It("should throw error", func() {
			client.ListBrokersReturns(&types.Brokers{Brokers: []types.Broker{broker}}, nil)
			err := executeWithArgs([]string{"non-existing", "{}"})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("broker with name non-existing not found"))
		})
	})

	Context("when list brokers returns an error", func() {
		It("should be handled", func() {
			client.ListBrokersReturns(nil, errors.New("error retrieving brokers"))
			err := executeWithArgs([]string{"non-existing", "{}"})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("error retrieving brokers"))
		})
	})

	Context("when name is not provided", func() {
		It("should throw error", func() {
			err := executeWithArgs([]string{})
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("[name] is required"))
		})
	})

	Context("when json is not provided", func() {
		It("should throw error", func() {
			err := executeWithArgs([]string{"broker"})
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("Nothing to update. Broker JSON is not provided"))
		})
	})

	Context("when json is invalid", func() {
		It("should throw error", func() {
			err := executeWithArgs([]string{"broker", "{name: none}"})
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("broker JSON is invalid"))
		})
	})

	Context("when output flag is used", func() {
		It("should print in json", func() {
			client.UpdateBrokerReturns(&broker, nil)
			client.ListBrokersReturns(&types.Brokers{Brokers: []types.Broker{broker}}, nil)
			executeWithArgs([]string{broker.Name, `{"name": "broker"}`, "-o", "json"})

			jsonByte, _ := json.MarshalIndent(broker, "", "  ")
			jsonOutputExpected := string(jsonByte) + "\n"
			Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
		})

		It("should print in yaml", func() {
			client.UpdateBrokerReturns(&broker, nil)
			client.ListBrokersReturns(&types.Brokers{Brokers: []types.Broker{broker}}, nil)
			executeWithArgs([]string{broker.Name, `{"name": "broker"}`, "-o", "yaml"})

			yamlByte, _ := yaml.Marshal(broker)
			yamlOutputExpected := string(yamlByte) + "\n"
			Expect(buffer.String()).To(ContainSubstring(yamlOutputExpected))
		})
	})
})
