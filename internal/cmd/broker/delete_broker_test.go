package broker

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

func TestDeleteBrokerCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Delete brokers command test", func() {

	var client *smclientfakes.FakeClient
	var command *DeleteBrokerCmd
	var buffer *bytes.Buffer
	var promptBuffer *bytes.Buffer

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		promptBuffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewDeleteBrokerCmd(context, promptBuffer)

		brokers := &types.Brokers{}
		brokers.Brokers = []types.Broker{{ID: "1234", Name: "broker-name"}, {ID: "456", Name: "broker2"}}
		client.ListBrokersWithQueryReturns(brokers, nil)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing broker is being deleted forcefully", func() {
		It("should list success message", func() {
			client.DeleteBrokerReturns(nil)
			err := executeWithArgs([]string{"broker-name", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Broker with name: broker-name successfully deleted"))
		})
	})

	Context("when existing broker is being deleted", func() {
		It("should list success message when confirmed", func() {
			client.DeleteBrokerReturns(nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs([]string{"broker-name"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Broker with name: broker-name successfully deleted"))
		})

		It("should print delete declined when declined", func() {
			client.DeleteBrokerReturns(nil)
			promptBuffer.WriteString("n")
			err := executeWithArgs([]string{"broker-name"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Delete declined"))
		})
	})

	Context("when 2 brokers are being deleted", func() {
		It("should print deleted ones", func() {
			client.DeleteBrokerReturns(nil)
			err := executeWithArgs([]string{"broker-name", "broker2", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Broker with name: broker-name successfully deleted"))
			Expect(buffer.String()).To(ContainSubstring("Broker with name: broker2 successfully deleted"))
		})
	})

	Context("when 2 brokers are being deleted and one is not found", func() {
		It("should print which was not found", func() {
			client.DeleteBrokerReturns(nil)
			err := executeWithArgs([]string{"broker-name", "broker", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Broker with name: broker-name successfully deleted"))
			Expect(buffer.String()).To(ContainSubstring("broker with name: broker was not found"))
		})
	})

	Context("when 3 brokers are being deleted and one is not found", func() {
		It("should print which was not found", func() {
			client.DeleteBrokerReturns(nil)
			err := executeWithArgs([]string{"broker-name", "non-existing", "broker2", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Broker with name: broker-name successfully deleted"))
			Expect(buffer.String()).To(ContainSubstring("Broker with name: broker2 successfully deleted"))
			Expect(buffer.String()).To(ContainSubstring("broker with name: non-existing was not found"))
		})
	})

	Context("when non-existing broker is being deleted", func() {
		It("should return error message", func() {
			expectedError := errors.ResponseError{StatusCode: http.StatusNotFound}
			client.ListBrokersWithQueryReturns(&types.Brokers{}, nil)
			client.DeleteBrokerReturns(expectedError)
			err := executeWithArgs([]string{"non-existing-name", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Broker(s) not found"))
		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeleteBrokerReturns(nil)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("[name] is required"))
		})
	})
})
