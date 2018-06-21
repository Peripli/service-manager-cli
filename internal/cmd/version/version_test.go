package version

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"fmt"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

func TestVersionCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Login Command test", func() {

	var command *Cmd
	var buffer *bytes.Buffer
	var clientVersion string

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		context := &cmd.Context{Output: buffer}
		clientVersion = "TEST VERSION"
		command = NewVersionCmd(context, clientVersion)
	})

	Describe("Valid request", func() {
		Context("Of the version command", func() {
			It("should print version", func() {
				vc := command.Prepare(cmd.CommonPrepare)
				err := vc.Execute()

				Expect(buffer.String()).To(Equal(fmt.Sprintf("Service Manager Client %s\n", clientVersion)))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
