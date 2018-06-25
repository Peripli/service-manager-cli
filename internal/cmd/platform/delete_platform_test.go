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

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewDeletePlatformCmd(context)

		platforms := &types.Platforms{}
		platforms.Platforms = []types.Platform{{ID: "1234", Name: "platform-name"}, {ID: "456", Name: "platform2"}}
		client.ListPlatformsReturns(platforms, nil)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing platform is being deleted", func() {
		It("should list success message", func() {
			client.DeletePlatformReturns(nil)
			err := executeWithArgs([]string{"platform-name"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Platform with name: platform-name successfully deleted"))
		})
	})

	Context("when 2 platforms are being deleted", func() {
		It("should success for both of them", func() {
			client.DeletePlatformReturns(nil)
			err := executeWithArgs([]string{"platform-name", "platform2"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Platform with name: platform-name successfully deleted"))
			Expect(buffer.String()).To(ContainSubstring("Platform with name: platform2 successfully deleted"))
		})
	})

	Context("when 2 platforms are being deleted and one is not found", func() {
		It("should print the name of the not found", func() {
			client.DeletePlatformReturns(nil)
			err := executeWithArgs([]string{"platform-name", "platform"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Platform with name: platform-name successfully deleted"))
			Expect(buffer.String()).To(ContainSubstring("Platform with name: platform was not found"))
		})
	})

	Context("when non-existing platform is being deleted", func() {
		It("should return error message", func() {
			expectedError := errors.ResponseError{StatusCode: http.StatusNotFound}
			client.DeletePlatformReturns(expectedError)
			err := executeWithArgs([]string{"non-existing-name"})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Platform(s) not found"))
		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeletePlatformReturns(nil)
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("[name] is required"))
		})
	})
})
