package login

import (
	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/auth/authfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"errors"
	"testing"

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
	var authStrategy *authfakes.FakeAuthenticator
	var client *smclientfakes.FakeClient

	authBuilder := func(options *auth.Options) (auth.Authenticator, *auth.Options, error) {
		return authStrategy, options, nil
	}

	BeforeEach(func() {
		client = &smclientfakes.FakeClient{}
		credentialsBuffer = &bytes.Buffer{}
		outputBuffer = &bytes.Buffer{}
		config = configurationfakes.FakeConfiguration{}
		authStrategy = &authfakes.FakeAuthenticator{}

		client.GetInfoReturns(&types.Info{TokenIssuerURL: "http://valid-uaa.com"}, nil)
		authStrategy.PasswordCredentialsReturns(&auth.Token{
			AccessToken: "access-token",
		}, nil)
		authStrategy.ClientCredentialsReturns(&auth.Token{
			AccessToken: "access-token",
		}, nil)

		context := &cmd.Context{Output: outputBuffer, Configuration: &config, Client: client}
		command = NewLoginCmd(context, credentialsBuffer, authBuilder)
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

				savedConfig := config.SaveArgsForCall(0)
				Expect(savedConfig.ClientID).To(Equal("cf"))
				Expect(savedConfig.ClientSecret).To(Equal(""))
			})
		})

		Context("With password and client id provided through flags", func() {
			It("should save configuration successfully", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				lc.SetArgs([]string{"--url=http://valid-url.com", "--password=password", "--client-id=smctl"})

				credentialsBuffer.WriteString("user\n")

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))

				savedConfig := config.SaveArgsForCall(0)
				Expect(savedConfig.ClientID).To(Equal("smctl"))
				Expect(savedConfig.ClientSecret).To(Equal(""))
			})
		})

		Context("With password, client id and client secret provided through flags", func() {
			It("should save configuration successfully", func() {
				lc := command.Prepare(cmd.CommonPrepare)
				lc.SetArgs([]string{"--url=http://valid-url.com", "--password=password", "--client-id=smctl", "--client-secret=smctl"})

				credentialsBuffer.WriteString("user\n")

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))

				savedConfig := config.SaveArgsForCall(0)
				Expect(savedConfig.ClientID).To(Equal("smctl"))
				Expect(savedConfig.ClientSecret).To(Equal("smctl"))
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

		Context("With client credentials flow", func() {
			When("client id secret are provided through flag", func() {
				It("login successfully", func() {
					lc := command.Prepare(cmd.CommonPrepare)
					lc.SetArgs([]string{"--url=http://valid-url.com", "--auth-flow=client-credentials", "--client-id=id", "--client-secret=secret"})

					err := lc.Execute()

					Expect(err).ShouldNot(HaveOccurred())
					Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))
				})
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

		Context("With client-credentials flow", func() {
			When("client id and secret is not provided", func() {
				It("should return an error", func() {
					lc := command.Prepare(cmd.CommonPrepare)
					lc.SetArgs([]string{"--url=http://valid-url.com", "--auth-flow=client-credentials"})

					err := lc.Execute()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(Equal("clientID/clientSecret should not be empty when using client credentials flow"))
				})
			})

			When("client id is not provided", func() {
				It("should return an error", func() {
					lc := command.Prepare(cmd.CommonPrepare)
					lc.SetArgs([]string{"--url=http://valid-url.com", "--auth-flow=client-credentials", "--client-secret", "secret"})

					err := lc.Execute()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(Equal("clientID/clientSecret should not be empty when using client credentials flow"))
				})
			})
		})
	})
})
