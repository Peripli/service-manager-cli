package broker

import (
	"github.com/Peripli/service-manager/pkg/util"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

var _ = Describe("Delete brokers command test", func() {
	var client *smclientfakes.FakeClient
	var command *DeleteBrokerCmd
	var buffer *bytes.Buffer
	var promptBuffer *bytes.Buffer
	var brokers *types.Brokers

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		promptBuffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewDeleteBrokerCmd(context, promptBuffer)

		brokers = &types.Brokers{}
		brokers.Brokers = []types.Broker{{ID: "1234", Name: "broker-name"}, {ID: "456", Name: "broker2"}}
	})

	JustBeforeEach(func() {
		client.ListBrokersReturns(brokers, nil)
	})

	executeWithArgs := func(args ...string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing broker is being deleted forcefully", func() {
		It("should list success message", func() {
			client.DeleteBrokerReturns("", nil)
			err := executeWithArgs("broker-name", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Broker successfully deleted."))
		})
	})

	Context("when existing broker is being deleted", func() {
		It("should list success message when confirmed", func() {
			client.DeleteBrokerReturns("", nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs("broker-name")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Broker successfully deleted."))
		})

		It("should print delete declined when declined", func() {
			client.DeleteBrokerReturns("", nil)
			promptBuffer.WriteString("n")
			err := executeWithArgs("broker-name")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Delete declined"))
		})
	})

	Context("when generic parameter flag is used", func() {
		It("should pass it to SM", func() {
			client.DeleteBrokerReturns("", nil)
			promptBuffer.WriteString("y")
			param := "parameterKey=parameterValue"
			err := executeWithArgs("broker-name", "--param", param)
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.DeleteBrokerArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param, "async=false"))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("With async flag", func() {
		It("should pass it to SM", func() {
			client.DeleteBrokerReturns("", nil)
			promptBuffer.WriteString("y")

			err := executeWithArgs("broker-name", "--mode", "async")
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.DeleteBrokerArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf("async=true"))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when non-existing brokers are being deleted", func() {
		BeforeEach(func() {
			brokers = &types.Brokers{}
		})
		It("should return message", func() {
			err := executeWithArgs("non-existing-name", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Broker not found"))
		})
	})

	Context("when SM returns error", func() {
		It("should return error message", func() {
			body := ioutil.NopCloser(bytes.NewReader([]byte("")))
			expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusInternalServerError})
			client.DeleteBrokerReturns("", expectedError)
			err := executeWithArgs("name", "-f")

			Expect(err).Should(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Could not delete broker. Reason:"))

		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeleteBrokerReturns("", nil)
			err := executeWithArgs()

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("single [name] is required"))
		})
	})
})
