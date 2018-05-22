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
)

func TestDeleteBrokerCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("List brokers command test", func() {

	var client *smclientfakes.FakeClient
	var command *DeleteBrokerCmd
	var buffer *bytes.Buffer

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewDeleteBrokerCmd(context)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Command()
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing broker is being deleted", func() {
		It("should list success message", func() {
			client.DeleteBrokerReturns(nil)
			err := executeWithArgs([]string{"broker-id"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Broker with id: broker-id successfully deleted"))
		})
	})

	Context("when non-existing broker is being deleted", func() {
		It("should return error message", func() {
			expectedError := errors.ResponseError{StatusCode: http.StatusNotFound}
			client.DeleteBrokerReturns(expectedError)
			err := executeWithArgs([]string{"invalid-id"})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeleteBrokerReturns(nil)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("[id] is required"))
		})
	})
})
