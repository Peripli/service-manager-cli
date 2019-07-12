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

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		context := &cmd.Context{Output: buffer}
		Version = "v1.2.3"
		GitCommit = "987654321"
		command = NewVersionCmd(context)
	})

	Describe("Valid request", func() {
		Context("Of the version command", func() {
			It("should print version", func() {
				vc := command.Prepare(cmd.CommonPrepare)
				err := vc.Execute()

				Expect(buffer.String()).To(Equal(fmt.Sprintf("Service Manager Client %s (%s)\n",
					Version, GitCommit)))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
