package login

import (
	"fmt"

	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/auth/authfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/spf13/cobra"

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
	var lc *cobra.Command
	var authOptions auth.Options

	authBuilder := func(options *auth.Options) (auth.Authenticator, *auth.Options, error) {
		authOptions = *options
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
		lc = command.Prepare(cmd.CommonPrepare)
	})

	Describe("Valid request", func() {
		Context("With password provided through flag", func() {
			It("should save configuration successfully with default client credentials", func() {
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
			It("should save configuration successfully without client credentials", func() {
				lc.SetArgs([]string{"--url=http://valid-url.com", "--password=password", "--client-id=smctl"})

				credentialsBuffer.WriteString("user\n")

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))

				savedConfig := config.SaveArgsForCall(0)
				Expect(savedConfig.ClientID).To(Equal(""))
				Expect(savedConfig.ClientSecret).To(Equal(""))
			})
		})

		Context("With password, client id and client secret provided through flags", func() {
			It("should save configuration successfully without client credentials", func() {
				lc.SetArgs([]string{"--url=http://valid-url.com", "--password=password", "--client-id=smctl", "--client-secret=smctl"})

				credentialsBuffer.WriteString("user\n")

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))

				savedConfig := config.SaveArgsForCall(0)
				Expect(savedConfig.ClientID).To(Equal(""))
				Expect(savedConfig.ClientSecret).To(Equal(""))
			})
		})

		Context("With user and password provided through flag", func() {
			It("should save configuration successfully", func() {
				lc.SetArgs([]string{"--url=http://valid-url.com", "--user=user", "--password=password"})

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))
			})
		})

		Context("With verbose flag provided", func() {
			It("should print more detailed messages", func() {
				lc.SetArgs([]string{"--url=http://valid-url.com", "--user=user", "--password=password"})

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))
			})
		})

		Context("With client credentials flow", func() {
			When("client id secret are provided through flag", func() {
				It("login successfully and not save the client credentials", func() {
					lc.SetArgs([]string{"--url=http://valid-url.com", "--auth-flow=client-credentials", "--client-id=id", "--client-secret=secret"})

					err := lc.Execute()

					Expect(err).ShouldNot(HaveOccurred())
					Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))

					savedConfig := config.SaveArgsForCall(0)
					Expect(savedConfig.ClientID).To(Equal(""))
					Expect(savedConfig.ClientSecret).To(Equal(""))
				})
			})
		})

		Context("Use token_basic_auth returned by info endpoint", func() {
			for _, tokenBasicAuth := range []bool{true, false} {
				tokenBasicAuth := tokenBasicAuth
				It(fmt.Sprintf("token_basic_auth: %v", tokenBasicAuth), func() {
					client.GetInfoReturns(&types.Info{
						TokenIssuerURL: "http://valid-uaa.com",
						TokenBasicAuth: tokenBasicAuth,
					},
						nil)
					lc.SetArgs([]string{"--url=http://valid-url.com", "--user=user", "--password=password"})

					err := lc.Execute()

					Expect(err).ShouldNot(HaveOccurred())
					savedConfig := config.SaveArgsForCall(0)
					Expect(authOptions.TokenBasicAuth).To(Equal(tokenBasicAuth))
					Expect(savedConfig.TokenBasicAuth).To(Equal(tokenBasicAuth))
				})
			}
		})

		FContext("With mTLS, cert & key & client id are provided through flags", func() {
			It("executes login and saved access token", func() {
				lc.SetArgs([]string{"--url=http://valid-url.com", "--auth-flow=client-credentials", "--cert=cert.pem", "--key=key.pem", "--client-id=smctl"})

				credentialsBuffer.WriteString("user\n")

				err := lc.Execute()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(outputBuffer.String()).To(ContainSubstring("Logged in successfully.\n"))

				savedConfig := config.SaveArgsForCall(0)
				Expect(savedConfig.ClientID).To(Equal(""))
				Expect(savedConfig.ClientSecret).To(Equal(""))
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With no URL flag provided", func() {
			It("should return error", func() {
				err := lc.Execute()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("URL flag must be provided"))
			})
		})

		Context("With invalid URL flag provided", func() {
			It("should return error", func() {
				lc.SetArgs([]string{"--url=htp://invalid-url.com"})
				err := lc.Execute()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("service manager URL is invalid"))
			})
		})

		Context("With empty username provided", func() {
			It("should return error", func() {
				lc.SetArgs([]string{"--url=http://valid-url.com", "--password=password"})
				credentialsBuffer.WriteString("\n")
				err := lc.Execute()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("username/password should not be empty"))
			})
		})

		Context("With invalid auth-flow provided", func() {
			It("should return error", func() {
				lc.SetArgs([]string{"--url=http://valid-url.com", "--auth-flow=bad-flow"})
				credentialsBuffer.WriteString("\n")
				err := lc.Execute()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("bad-flow"))
			})
		})

		Context("With error while typing user in", func() {
			It("should save configuration successfully", func() {

				err := lc.Execute()

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("With error while saving configuration", func() {
			It("should return error", func() {
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
					lc.SetArgs([]string{"--url=http://valid-url.com", "--auth-flow=client-credentials"})

					err := lc.Execute()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(Equal("clientID/clientSecret should not be empty when using client credentials flow"))
				})
			})

			When("client id is not provided", func() {
				It("should return an error", func() {
					lc.SetArgs([]string{"--url=http://valid-url.com", "--auth-flow=client-credentials", "--client-secret", "secret"})

					err := lc.Execute()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(Equal("clientID/clientSecret should not be empty when using client credentials flow"))
				})
			})
		})
	})
})
