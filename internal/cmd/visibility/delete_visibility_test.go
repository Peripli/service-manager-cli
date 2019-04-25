package visibility

import (
	"bytes"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

func TestDeleteVisibilityCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Delete visibility command test", func() {

	var client *smclientfakes.FakeClient
	var command *DeleteVisibilityCmd
	var buffer *bytes.Buffer
	var promptBuffer *bytes.Buffer

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		promptBuffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewDeleteVisibilityCmd(context, promptBuffer)
	})

	executeWithArgs := func(args... string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing visibility is being deleted forcefully", func() {
		It("should list success message", func() {
			client.DeleteVisibilityReturns(nil)
			err := executeWithArgs("id", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Visibility with id: id successfully deleted"))
		})
	})

	Context("when existing visibility is being deleted", func() {
		It("should list success message when confirmed", func() {
			client.DeleteVisibilityReturns(nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs("id")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Visibility with id: id successfully deleted"))
		})

		It("should print delete declined when declined", func() {
			client.DeleteVisibilityReturns(nil)
			promptBuffer.WriteString("n")
			err := executeWithArgs("id")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Delete declined"))
		})
	})

	Context("when 2 visibilities are being deleted", func() {
		It("should success for both of them", func() {
			client.DeleteVisibilityReturns(nil)
			err := executeWithArgs("id1", "id2", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Visibility with id: id1 successfully deleted"))
			Expect(buffer.String()).To(ContainSubstring("Visibility with id: id2 successfully deleted"))
		})
	})

	Context("when non-existing visibility is being deleted", func() {
		It("should return error message", func() {
			expectedError := errors.ResponseError{StatusCode: http.StatusNotFound}
			client.DeleteVisibilityReturns(expectedError)
			err := executeWithArgs("id", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("visibility with id: id was not found"))
		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeleteVisibilityReturns(nil)
			err := executeWithArgs()

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("id is required"))
		})
	})
})
