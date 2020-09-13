/*
 * Copyright 2018 The Service Manager Authors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package login

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"syscall"

	cliErr "github.com/Peripli/service-manager-cli/pkg/errors"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/configuration"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/internal/util"
	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	defaultClientID     = "cf"
	defaultClientSecret = ""
)

// Cmd wraps the smctl login command
type Cmd struct {
	*cmd.Context

	input io.ReadWriter

	serviceManagerURL  string
	user               string
	password           string
	sslDisabled        bool
	clientID           string
	clientSecret       string
	authenticationFlow auth.Flow

	authBuilder authenticationBuilder
}

type authenticationBuilder func(*auth.Options) (auth.Authenticator, *auth.Options, error)

// NewLoginCmd return new login command with context and input reader
func NewLoginCmd(context *cmd.Context, input io.ReadWriter, authBuilder authenticationBuilder) *Cmd {
	return &Cmd{Context: context, input: input, authBuilder: authBuilder}
}

// Prepare returns cobra command
func (lc *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "login",
		Aliases: []string{"l"},
		Short:   "Logs user in",
		Long:    `Connects to a Service Manager and logs user in.`,

		PreRunE: prepare(lc, lc.Context),
		RunE:    cmd.RunE(lc),
	}

	result.Flags().StringVarP(&lc.serviceManagerURL, "url", "a", "", "Base URL of the Service Manager")
	result.Flags().StringVarP(&lc.user, "user", "u", "", "User ID")
	result.Flags().StringVarP(&lc.password, "password", "p", "", "Password")
	result.Flags().StringVarP(&lc.clientID, "client-id", "", "", "Client id used for OAuth flow")
	result.Flags().StringVarP(&lc.clientSecret, "client-secret", "", defaultClientSecret, "Client secret used for OAuth flow")
	result.Flags().BoolVarP(&lc.sslDisabled, "skip-ssl-validation", "", false, "Skip verification of the OAuth endpoint. Not recommended!")
	result.Flags().StringVarP((*string)(&lc.authenticationFlow), "auth-flow", "", string(auth.PasswordGrant), `Authentication flow (grant type): "client-credentials" or "password-grant"`)
	cmd.AddCommonQueryFlag(result.Flags(), &lc.Parameters)

	return result
}

// HideUsage hides the command's usage
func (lc *Cmd) HideUsage() bool {
	return true
}

// Validate valides the command's arguments
func (lc *Cmd) Validate(args []string) error {
	if lc.serviceManagerURL == "" {
		return errors.New("URL flag must be provided")
	}

	if err := util.ValidateURL(lc.serviceManagerURL); err != nil {
		return fmt.Errorf("service manager URL is invalid: %v", err)
	}

	if err := lc.validateLoginFlow(); err != nil {
		return err
	}

	return nil
}

// Run runs the logic of the command
func (lc *Cmd) Run() error {
	httpClient := util.BuildHTTPClient(lc.sslDisabled)

	if lc.Client == nil {
		lc.Client = smclient.NewClient(lc.Ctx, httpClient, lc.serviceManagerURL)
	}

	info, err := lc.Client.GetInfo(&lc.Parameters)
	if err != nil {
		return cliErr.New("Could not get Service Manager info", err)
	}

	options := &auth.Options{
		User:           lc.user,
		Password:       lc.password,
		ClientID:       lc.clientID,
		ClientSecret:   lc.clientSecret,
		IssuerURL:      info.TokenIssuerURL,
		TokenBasicAuth: info.TokenBasicAuth,
		SSLDisabled:    lc.sslDisabled,
	}

	authStrategy, options, err := lc.authBuilder(options)
	if err != nil {
		return cliErr.New("Could not build authenticator", err)
	}
	token, err := lc.getToken(authStrategy)
	if err != nil {
		description := "could not login"
		if len(lc.Parameters.GeneralParams) == 0 {
			description += " in case of tenant subdomain param required"
		}
		return cliErr.New(description, err)
	}

	settings := &configuration.Settings{
		URL:         lc.serviceManagerURL,
		User:        lc.user,
		SSLDisabled: lc.sslDisabled,
		AuthFlow:    lc.authenticationFlow,

		Token: *token,

		IssuerURL:             info.TokenIssuerURL,
		AuthorizationEndpoint: options.AuthorizationEndpoint,
		TokenEndpoint:         options.TokenEndpoint,
		TokenBasicAuth:        info.TokenBasicAuth,
	}
	if options.ClientID == defaultClientID && options.ClientSecret == defaultClientSecret {
		settings.ClientID = options.ClientID
		settings.ClientSecret = options.ClientSecret
	}
	if settings.User == "" {
		settings.User = options.ClientID
	}
	err = lc.Configuration.Save(settings)

	if err != nil {
		return err
	}

	output.PrintMessage(lc.Output, "Logged in successfully.\n")
	return nil
}

func (lc *Cmd) getToken(authStrategy auth.Authenticator) (*auth.Token, error) {
	switch lc.authenticationFlow {
	case auth.ClientCredentials:
		return authStrategy.ClientCredentials()
	case auth.PasswordGrant:
		return authStrategy.PasswordCredentials(lc.user, lc.password)
	default:
		return nil, fmt.Errorf("authentication flow %s not recognized", lc.authenticationFlow)
	}
}

func (lc *Cmd) validateLoginFlow() error {
	switch lc.authenticationFlow {
	case auth.ClientCredentials:
		if len(lc.clientID) == 0 || len(lc.clientSecret) == 0 {
			return errors.New("clientID/clientSecret should not be empty when using client credentials flow")
		}
	case auth.PasswordGrant:
		return lc.validatePasswordGrant()
	default:
		return fmt.Errorf("unknown authentication flow: %s", lc.authenticationFlow)
	}

	return nil
}

func (lc *Cmd) validatePasswordGrant() error {
	if len(lc.clientID) == 0 {
		lc.clientID = defaultClientID
	}

	if err := lc.readUser(); err != nil {
		return err
	}

	if err := lc.readPassword(); err != nil {
		return err
	}

	if len(lc.user) == 0 || len(lc.password) == 0 {
		return errors.New("username/password should not be empty")
	}
	return nil
}

func (lc *Cmd) readUser() error {
	if lc.user == "" {
		output.PrintMessage(lc.Output, "User: ")
		bufReader := bufio.NewReader(lc.input)
		readUser, isPrefix, err := bufReader.ReadLine()
		if isPrefix {
			return errors.New("username too long")
		}
		if err != nil {
			return err
		}

		lc.user = string(readUser)
	}
	return nil
}

func (lc *Cmd) readPassword() error {
	if lc.password == "" {
		output.PrintMessage(lc.Output, "Password: ")
		password, err := terminal.ReadPassword((int)(syscall.Stdin))
		output.Println(lc.Output)
		if err != nil {
			return err
		}

		lc.password = string(password)
	}
	return nil
}
