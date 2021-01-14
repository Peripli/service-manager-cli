package logout

import (
	"fmt"
	"github.com/Peripli/service-manager-cli/pkg/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	"bytes"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/configuration"
	"github.com/Peripli/service-manager-cli/internal/configuration/configurationfakes"
)

func TestInfoCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Logout Command test", func() {

	var command *Cmd
	var buffer *bytes.Buffer
	var config *configurationfakes.FakeConfiguration

	settings := &configuration.Settings{
		Token: auth.Token{
			AccessToken:  "access-token",
			TokenType:    "token-type",
			RefreshToken: "token-refresh",
			ExpiresIn:    time.Time{},
			Scope:        "",
		},
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		config = &configurationfakes.FakeConfiguration{}
		context := &cmd.Context{Output: buffer, Configuration: config}
		command = NewLogoutCmd(context)
	})

	Describe("user logout", func() {
		Context("With logged user", func() {
			It("should delete token details from settings", func() {
				config.LoadReturns(settings, nil)

				ic := command.Prepare(cmd.CommonPrepare)
				err := ic.Execute()
				Expect(err).ShouldNot(HaveOccurred())

				config.LoadReturns(settings, nil)

				Expect(settings.Scope).To(Equal(""))
				Expect(settings.AccessToken).To(Equal(""))
				Expect(settings.RefreshToken).To(Equal(""))
				Expect(settings.ExpiresIn).To(Equal(time.Time{}))
			})
		})

		Context("with no logged in user", func() {
			BeforeEach(func() {
				settings = &configuration.Settings{
					Token: auth.Token{},
				}
				err := config.Save(settings)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should indicate that no user is logged in and do nothing", func() {
				config.LoadReturns(settings, fmt.Errorf("error"))
				ic := command.Prepare(cmd.CommonPrepare)
				err := ic.Execute()
				Expect(err).ShouldNot(HaveOccurred())
				print(buffer.String())
				Expect(buffer.String()).To(ContainSubstring("You are already logged out."))
			})
		})
	})
})
