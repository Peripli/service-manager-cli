package platform

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

func TestDeletePlatformCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Delete platforms command test", func() {

	var client *smclientfakes.FakeClient
	var command *DeletePlatformCmd
	var buffer *bytes.Buffer
	var promptBuffer *bytes.Buffer

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		promptBuffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewDeletePlatformCmd(context, promptBuffer)

		platforms := &types.Platforms{}
		platforms.Platforms = []types.Platform{{ID: "1234", Name: "platform-name"}, {ID: "456", Name: "platform2"}}
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing platform is being deleted forcefully", func() {
		It("should list success message", func() {
			client.DeletePlatformsByFieldQueryReturns(nil)
			err := executeWithArgs([]string{"platform-name", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Platform(s) successfully deleted."))
		})
	})

	Context("when existing platform is being deleted", func() {
		It("should list success message when confirmed", func() {
			client.DeletePlatformsByFieldQueryReturns(nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs([]string{"platform-name"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Platform(s) successfully deleted."))
		})

		It("should print delete declined when declined", func() {
			client.DeletePlatformsByFieldQueryReturns(nil)
			promptBuffer.WriteString("n")
			err := executeWithArgs([]string{"platform-name"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Delete declined"))
		})

	})

	Context("when non-existing platform is being deleted", func() {
		It("should return error message", func() {
			expectedError := errors.ResponseError{StatusCode: http.StatusNotFound}
			client.ListPlatformsWithQueryReturns(&types.Platforms{}, nil)
			client.DeletePlatformsByFieldQueryReturns(expectedError)
			err := executeWithArgs([]string{"non-existing-name", "-f"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Platform(s) not found."))
		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeletePlatformsByFieldQueryReturns(nil)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("[name] is required"))
		})
	})
})
