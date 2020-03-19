package broker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

var _ = Describe("Get broker command test", func() {

	var client *smclientfakes.FakeClient
	var command *GetBrokerCmd
	var buffer *bytes.Buffer
	broker := types.Broker{
		Name: "broker1",
		ID:   "id1",
		URL:  "http://broker1.com",
	}
	broker2 := types.Broker{
		Name: "broker2",
		ID:   "id2",
		URL:  "http://broker2.com",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		client.ListBrokersReturns(&types.Brokers{Brokers: []types.Broker{broker, broker2}}, nil)
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewGetBrokerCmd(context)
	})

	executeWithArgs := func(args ...string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no broker name is provided", func() {
		It("should return error", func() {
			client.GetBrokerByIDReturns(&broker, nil)
			err := executeWithArgs("")

			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when no known broker is provided", func() {
		It("should return no brokers", func() {
			client.ListBrokersReturns(&types.Brokers{}, nil)
			err := executeWithArgs("unknown")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("No broker found with name: unknown"))
		})
	})

	Context("when broker with name is found", func() {
		It("should return its data", func() {
			client.GetBrokerByIDReturns(&broker, nil)
			err := executeWithArgs("broker1")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(broker.TableData().String()))
		})
	})
})
