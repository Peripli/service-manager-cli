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
	"encoding/base64"
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
}

// NewLoginCmd return new login command with context and input reader
func NewLoginCmd(context *cmd.Context, input io.ReadWriter) *Cmd {
	return &Cmd{Context: context, input: input}
}

// HideUsage hides the command's usage
func (lc *Cmd) HideUsage() bool {
	return true
}

// Command returns cobra command
func (lc *Cmd) Command() *cobra.Command {
	result := lc.buildCommand()
	result = lc.addFlags(result)

	return result
}

func (lc *Cmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "login",
		Aliases: []string{"l"},
		Short:   "Logs user in",
		Long:    `Connects to a Service Manager and logs user in.`,

		PreRunE: cmd.PreRunE(lc, lc.Context),
		RunE:    cmd.RunE(lc),
	}
}

func (lc *Cmd) addFlags(command *cobra.Command) *cobra.Command {
	command.PersistentFlags().StringVarP(&lc.serviceManagerURL, "url", "a", "", "Base URL of the Service Manager")
	command.PersistentFlags().StringVarP(&lc.user, "user", "u", "", "User ID")
	command.PersistentFlags().StringVarP(&lc.password, "password", "p", "", "Password")

	return command
}

// Validate valides the command's arguments
func (lc *Cmd) Validate(args []string) error {
	if lc.serviceManagerURL == "" {
		return errors.New("URL flag must be provided")
	}
	return nil
}

// Run runs the logic of the command
func (lc *Cmd) Run() error {
	if err := util.ValidateURL(lc.serviceManagerURL); err != nil {
		return fmt.Errorf("service manager URL is invalid: %v", err)
	}

	if err := lc.readUser(); err != nil {
		return err
	}

	if err := lc.readPassword(); err != nil {
		return err
	}

	if lc.Verbose {
		output.PrintMessage(lc.Output, "Connecting to Service Manager: %s\n", lc.serviceManagerURL)
	}

	token := "basic " + base64.StdEncoding.EncodeToString([]byte(lc.user+":"+lc.password))
	err := lc.Configuration.Save(&smclient.ClientConfig{URL: lc.serviceManagerURL, User: lc.user, Token: token})
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

		lc.user = (string)(readUser)
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
