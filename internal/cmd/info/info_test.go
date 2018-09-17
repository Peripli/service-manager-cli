package info

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"fmt"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/configuration/configurationfakes"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
)

func TestInfoCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Login Command test", func() {

	var command *Cmd
	var buffer *bytes.Buffer
	var config *configurationfakes.FakeConfiguration

	var clientConfig smclient.ClientConfig

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		config = &configurationfakes.FakeConfiguration{}
		context := &cmd.Context{Output: buffer, Configuration: config}
		command = NewInfoCmd(context)
		clientConfig = smclient.ClientConfig{URL: "http://test-url.com", User: "test-user"}
	})

	Describe("Valid request", func() {

		JustBeforeEach(func() {
			config.UnmarshalKeyReturns(nil)
			config.UnmarshalKeyStub = func(key string, value interface{}) error {
				val, ok := value.(*smclient.ClientConfig)
				if ok {
					*val = clientConfig
				}
				return nil
			}
		})

		Context("With no logged user", func() {
			It("should print prompt to log in", func() {
				clientConfig.User = ""

				ic := command.Prepare(cmd.CommonPrepare)
				err := ic.Execute()

				Expect(buffer.String()).To(Equal("There is no logged user. Use \"smctl login\" to log in.\n"))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("With logged user", func() {
			It("should print URL and logged user", func() {
				ic := command.Prepare(cmd.CommonPrepare)
				err := ic.Execute()

				Expect(buffer.String()).To(ContainSubstring(fmt.Sprintf("Service Manager URL: %s\n", clientConfig.URL)))
				Expect(buffer.String()).To(ContainSubstring(fmt.Sprintf("Logged user: %s\n", clientConfig.User)))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
