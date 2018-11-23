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

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/internal/util"
	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/Peripli/service-manager/pkg/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	defaultClientID     = "cf"
	defaultClientSecret = ""
)

type authFlow string

func (a *authFlow) String() string {
	return string(*a)
}
func (a *authFlow) Set(value string) error {
	*a = authFlow(value)
	return nil
}

func (a *authFlow) Type() string {
	return "authFlow"
}

func newAuthFlowValue(value authFlow, p *authFlow) *authFlow {
	*p = value
	return (*authFlow)(p)
}

const (
	clientCredentials authFlow = "client-credentials"
	passwordGrant     authFlow = "password-grant"
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
	authenticationFlow authFlow

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
	result.Flags().VarP(newAuthFlowValue(passwordGrant, &lc.authenticationFlow), "auth-flow", "", "provide Oauth2 authentication flow type")

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

	return nil
}

// Run runs the logic of the command
func (lc *Cmd) Run() error {
	httpClient := util.BuildHTTPClient(lc.sslDisabled)
	clientConfig := &smclient.ClientConfig{
		URL: lc.serviceManagerURL,
	}
	if lc.Client == nil {
		lc.Client = smclient.NewClient(httpClient, clientConfig)
	}

	info, err := lc.Client.GetInfo()
	if err != nil {
		return cmd.NewError("Could not get Service Manager info", err.Error())
	}

	if err := lc.checkLoginFlow(); err != nil {
		return err
	}

	options := &auth.Options{
		ClientID:     lc.clientID,
		ClientSecret: lc.clientSecret,
		IssuerURL:    info.TokenIssuerURL,
		SSLDisabled:  lc.sslDisabled,
	}

	authStrategy, options, err := lc.authBuilder(options)
	if err != nil {
		return cmd.NewError("Could not build authenticator", err.Error())
	}
	token, err := lc.getToken(authStrategy)
	if err != nil {
		if authError, ok := err.(*auth.Error); ok {
			log.D().Debugf("login: %s", authError.Error())
			return cmd.NewError("Could not login", authError.Description)
		}
		return cmd.NewError("Could not login", err.Error())
	}

	clientConfig = &smclient.ClientConfig{
		URL:         lc.serviceManagerURL,
		User:        lc.user,
		SSLDisabled: lc.sslDisabled,

		Token:        *token,
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,

		IssuerURL:             info.TokenIssuerURL,
		AuthorizationEndpoint: options.AuthorizationEndpoint,
		TokenEndpoint:         options.TokenEndpoint,
	}
	if clientConfig.User == "" {
		clientConfig.User = options.ClientID
	}
	err = lc.Configuration.Save(clientConfig)

	if err != nil {
		return err
	}

	output.PrintMessage(lc.Output, "Logged in successfully.\n")
	return nil
}

func (lc *Cmd) getToken(authStrategy auth.Authenticator) (*auth.Token, error) {
	switch lc.authenticationFlow {
	case clientCredentials:
		return authStrategy.ClientCredentials()
	case passwordGrant:
		return authStrategy.PasswordCredentials(lc.user, lc.password)
	default:
		return nil, fmt.Errorf("authentication flow %s not recognized", lc.authenticationFlow)
	}
}

func (lc *Cmd) checkLoginFlow() error {
	if lc.authenticationFlow == clientCredentials {
		if len(lc.clientID) == 0 || len(lc.clientSecret) == 0 {
			return errors.New("clientID/clientSecret should not be empty when using client credentials flow")
		}
	} else {
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
