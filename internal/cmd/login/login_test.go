package login

import (
	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"errors"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/auth"
	"github.com/Peripli/service-manager-cli/internal/auth/authfakes"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/configuration/configurationfakes"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
)

func TestLoginCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Login Command test", func() {

	var command *Cmd
	var credentialsBuffer, outputBuffer *bytes.Buffer
	var config configurationfakes.FakeConfiguration
	var authStrategy authfakes.FakeAuthenticationStrategy
	var client *smclientfakes.FakeClient

	BeforeEach(func() {
		client = &smclientfakes.FakeClient{}
		credentialsBuffer = &bytes.Buffer{}
		outputBuffer = &bytes.Buffer{}
		config = configurationfakes.FakeConfiguration{}
		authStrategy = authfakes.FakeAuthenticationStrategy{}

		client.GetInfoReturns(&types.Info{TokenIssuerURL: "http://valid-uaa.com"}, nil)
		authStrategy.AuthenticateReturns(&auth.Token{AccessToken: "access-token"}, nil)

		context := &cmd.Context{Output: outputBuffer, Configuration: &config, Client: client}
		command = NewLoginCmd(context, credentialsBuffer, &authStrategy)
	})

	Describe("Valid request", func() {
		Context("With password provided through flag", func() {
			It("should save configuration successfully", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				lc.SetArgs([]string{"--url=http://valid-url.com", "--password=password"})

				credentialsBuffer.WriteString("user\n")

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))
			})
		})

		Context("With user and password provided through flag", func() {
			It("should save configuration successfully", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				lc.SetArgs([]string{"--url=http://valid-url.com", "--user=user", "--password=password"})

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))
			})
		})

		Context("With verbose flag provided", func() {
			It("should print more detailed messages", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				lc.SetArgs([]string{"--url=http://valid-url.com", "--user=user", "--password=password"})

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With no URL flag provided", func() {
			It("should return error", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				err := lc.Execute()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("URL flag must be provided"))
			})
		})

		Context("With invalid URL flag provided", func() {
			It("should return error", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				lc.SetArgs([]string{"--url=htp://invalid-url.com"})
				err := lc.Execute()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("service manager URL is invalid"))
			})
		})

		Context("With empty username provided", func() {
			It("should return error", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				lc.SetArgs([]string{"--url=http://valid-url.com", "--password=password"})
				credentialsBuffer.WriteString("\n")
				err := lc.Execute()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("username/password should not be empty"))
			})
		})

		Context("With error while typing user in", func() {
			It("should save configuration successfully", func() {
				lc := command.Prepare(cmd.CommonPrepare)

				err := lc.Execute()

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("With error while saving configuration", func() {
			It("should return error", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				lc.SetArgs([]string{"--url=http://valid-url.com", "--user=user", "--password=password"})
				config.SaveReturns(errors.New("saving configuration error"))

				err := lc.Execute()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("saving configuration error"))
			})
		})
	})
})
