package info

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/configuration"
	"github.com/Peripli/service-manager-cli/internal/configuration/configurationfakes"
)

func TestInfoCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Login Command test", func() {

	var command *Cmd
	var buffer *bytes.Buffer
	var config *configurationfakes.FakeConfiguration

	settings := &configuration.Settings{URL: "http://test-url.com", User: "test-user"}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		config = &configurationfakes.FakeConfiguration{}
		context := &cmd.Context{Output: buffer, Configuration: config}
		command = NewInfoCmd(context)
	})

	Describe("Valid request", func() {
		Context("With no logged user", func() {
			It("should print prompt to log in", func() {
				config.LoadReturns(nil, errors.New("configuration file not found"))

				ic := command.Prepare(cmd.CommonPrepare)
				err := ic.Execute()

				Expect(buffer.String()).To(Equal("There is no logged user. Use \"smctl login\" to log in.\n"))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("With logged user", func() {
			It("should print URL and logged user", func() {
				config.LoadReturns(settings, nil)

				ic := command.Prepare(cmd.CommonPrepare)
				err := ic.Execute()

				Expect(buffer.String()).To(ContainSubstring(fmt.Sprintf("Service Manager URL: %s\n", settings.URL)))
				Expect(buffer.String()).To(ContainSubstring(fmt.Sprintf("Logged user: %s\n", settings.User)))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
