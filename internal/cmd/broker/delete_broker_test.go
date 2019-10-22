package broker

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Peripli/service-manager/pkg/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
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
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing broker is being deleted forcefully", func() {
		It("should list success message", func() {
			client.DeleteBrokersReturns(nil)
			err := executeWithArgs([]string{"broker-name", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Broker(s) successfully deleted."))
		})
	})

	Context("when existing broker is being deleted", func() {
		It("should list success message when confirmed", func() {
			client.DeleteBrokersReturns(nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs([]string{"broker-name"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Broker(s) successfully deleted."))
		})

		It("should print delete declined when declined", func() {
			client.DeleteBrokersReturns(nil)
			promptBuffer.WriteString("n")
			err := executeWithArgs([]string{"broker-name"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Delete declined"))
		})
	})

	Context("when generic parameter flag is used", func() {
		It("should pass it to SM", func() {
			client.DeleteBrokersReturns(nil)
			promptBuffer.WriteString("y")
			param := "parameterKey=parameterValue"
			err := executeWithArgs([]string{"broker-name", "--param", param})
			Expect(err).ShouldNot(HaveOccurred())

			args := client.DeleteBrokersArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param))
			Expect(args.FieldQuery).To(ConsistOf("name eq 'broker-name'"))
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when non-existing brokers are being deleted", func() {
		It("should return message", func() {
			body := ioutil.NopCloser(bytes.NewReader([]byte("")))
			expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusNotFound})
			client.DeleteBrokersReturns(expectedError)
			err := executeWithArgs([]string{"non-existing-name", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Broker(s) not found"))
		})
	})

	Context("when SM returns error", func() {
		It("should return error message", func() {
			body := ioutil.NopCloser(bytes.NewReader([]byte("")))
			expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusInternalServerError})
			client.DeleteBrokersReturns(expectedError)
			err := executeWithArgs([]string{"name", "-f"})

			Expect(err).Should(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Could not delete broker(s). Reason:"))

		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeleteBrokersReturns(nil)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("single [name] is required"))
		})
	})
})
