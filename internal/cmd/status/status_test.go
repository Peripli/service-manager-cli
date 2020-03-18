package status

import (
	"github.com/Peripli/service-manager/pkg/util"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

func TestStatusCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Status command test", func() {

	var client *smclientfakes.FakeClient
	var command *Cmd
	var buffer *bytes.Buffer
	operation := &types.Operation{
		ID:         "operation-id",
		Type:       "create",
		State:      "failed",
		ResourceID: "broker-id",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewStatusCmd(context)
	})

	executeWithArgs := func(args ...string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no operation url is provided", func() {
		It("should return error", func() {
			client.StatusReturns(operation, nil)
			err := executeWithArgs()

			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when operation is not found", func() {
		It("should return message", func() {
			body := ioutil.NopCloser(bytes.NewReader([]byte("")))
			expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusNotFound})
			client.StatusReturns(nil, expectedError)
			err := executeWithArgs("non-existing-path")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Operation not found"))
		})
	})

	Context("when operation is found", func() {
		It("should return its data", func() {
			client.StatusReturns(operation, nil)
			err := executeWithArgs("path")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(operation.TableData().String()))
		})
	})
})
