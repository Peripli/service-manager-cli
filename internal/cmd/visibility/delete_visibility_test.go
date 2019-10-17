package visibility

import (
	"bytes"

	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	executeWithArgs := func(args ...string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing visibility is being deleted forcefully", func() {
		It("should list success message", func() {
			client.DeleteVisibilitiesReturns(nil)
			err := executeWithArgs("id", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Visibility successfully deleted."))
		})
	})

	Context("when existing visibility is being deleted", func() {
		It("should list success message when confirmed", func() {
			client.DeleteVisibilitiesReturns(nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs("id")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Visibility successfully deleted."))
		})

		It("should print delete declined when declined", func() {
			client.DeleteVisibilitiesReturns(nil)
			promptBuffer.WriteString("n")
			err := executeWithArgs("id")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Delete declined"))
		})
	})

	Context("when generic parameter flag is used", func() {
		It("should pass it to SM", func() {
			client.DeleteVisibilitiesReturns(nil)
			promptBuffer.WriteString("y")
			param := "parameterKey=parameterValue"
			err := executeWithArgs([]string{"id", "-f", "--param", param}...)
			Expect(err).ShouldNot(HaveOccurred())

			args := client.DeleteVisibilitiesArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param))
			Expect(args.FieldQuery).To(ConsistOf("id = id"))
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when non-existing visibility is being deleted", func() {
		It("should return error message", func() {
			body := ioutil.NopCloser(bytes.NewReader([]byte("")))
			expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusNotFound})
			client.DeleteVisibilitiesReturns(expectedError)
			err := executeWithArgs("id", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Visibility not found."))
		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeleteVisibilitiesReturns(nil)
			err := executeWithArgs()

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("single [id] is required"))
		})
	})
})
