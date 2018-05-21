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
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"syscall"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/print"
	"github.com/Peripli/service-manager-cli/internal/util"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// Wraps the smctl login command
type LoginCmd struct {
	*cmd.Context

	input io.ReadWriter

	serviceManagerURL string
	user              string
	password          string
}

func NewLoginCmd(context *cmd.Context, input io.ReadWriter) *LoginCmd {
	return &LoginCmd{Context: context, input: input}
}

func (rpc *LoginCmd) HideUsage() bool {
	return true
}

func (lc *LoginCmd) Command() *cobra.Command {
	result := lc.buildCommand()
	result = lc.addFlags(result)

	return result
}

func (lc *LoginCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "login",
		Aliases: []string{"l"},
		Short:   "Logs user in",
		Long:    `Connects to a Service Manager and logs user in.`,

		PreRunE: cmd.PreRunE(lc, lc.Context),
		RunE:    cmd.RunE(lc),
	}
}

func (lc *LoginCmd) addFlags(command *cobra.Command) *cobra.Command {
	command.PersistentFlags().StringVarP(&lc.serviceManagerURL, "url", "a", "", "Base URL of the Service Manager")
	command.PersistentFlags().StringVarP(&lc.user, "user", "u", "", "User ID")
	command.PersistentFlags().StringVarP(&lc.password, "password", "p", "", "Password")

	return command
}

func (lc *LoginCmd) Validate(args []string) error {
	if lc.serviceManagerURL == "" {
		return errors.New("URL flag must be provided")
	}
	return nil
}

func (lc *LoginCmd) Run() error {
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
		print.PrintMessage(lc.Output, "Connecting to Service Manager: %s\n", lc.serviceManagerURL)
	}

	token := "basic " + base64.StdEncoding.EncodeToString([]byte(lc.user+":"+lc.password))
	err := lc.Configuration.Save(&smclient.ClientConfig{lc.serviceManagerURL, lc.user, token})
	if err != nil {
		return err
	}

	print.PrintMessage(lc.Output, "Logged in successfully.\n")
	return nil
}

func (lc *LoginCmd) readUser() error {
	if lc.user == "" {
		print.PrintMessage(lc.Output, "User: ")
		readUser := make([]byte, 256)
		n, err := lc.input.Read(readUser)
		if err != nil {
			return err
		}

		lc.user = (string)(readUser[:n-1])
	}
	return nil
}

func (lc *LoginCmd) readPassword() error {
	if lc.password == "" {
		print.PrintMessage(lc.Output, "Password: ")
		password, err := terminal.ReadPassword((int)(syscall.Stdin))
		print.Println(lc.Output)
		if err != nil {
			return err
		}

		lc.password = string(password)
	}
	return nil
}
