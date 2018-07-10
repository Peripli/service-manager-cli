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
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// Cmd wraps the smctl login command
type Cmd struct {
	*cmd.Context

	input io.ReadWriter

	serviceManagerURL string
	user              string
	password          string
	sslDisabled		  bool
}

// NewLoginCmd return new login command with context and input reader
func NewLoginCmd(context *cmd.Context, input io.ReadWriter) *Cmd {
	return &Cmd{Context: context, input: input}
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
	result.Flags().BoolVarP(&lc.sslDisabled, "skip-ssl-validation", "", false, "Skip verification of the OAuth endpoint. Not recommended!")

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
	if lc.Client == nil {
		clientConfig := &smclient.ClientConfig{URL: lc.serviceManagerURL}
		lc.Client = smclient.NewClient(lc.Ctx, clientConfig)
	}

	info, err := lc.Client.GetInfo()
	if err != nil {
		return err
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

	config, token, err := lc.AuthStrategy.Authenticate(lc.Ctx, info.TokenIssuerURL, lc.user, lc.password)
	if err != nil {
		return err
	}

	err = lc.Configuration.Save(&smclient.ClientConfig{
		URL:    lc.serviceManagerURL,
		User:   lc.user,
		SSLDisabled: lc.sslDisabled,
		Token:  *token,
		Config: *config,
	})

	if err != nil {
		return err
	}

	output.PrintMessage(lc.Output, "Logged in successfully.\n")
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
