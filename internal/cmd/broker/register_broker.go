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

package broker

import (
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"fmt"
	"strings"

	"errors"

	"github.com/spf13/cobra"
)

// RegisterBrokerCmd wraps the smctl register-broker command
type RegisterBrokerCmd struct {
	*cmd.Context

	broker types.Broker

	basicString  string
	outputFormat int
}

// NewRegisterBrokerCmd returns new register-broker command with context
func NewRegisterBrokerCmd(context *cmd.Context) *RegisterBrokerCmd {
	return &RegisterBrokerCmd{Context: context, broker: types.Broker{}}
}

// Prepare returns cobra command
func (rbc *RegisterBrokerCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "register-broker [name] [url] <description>",
		Aliases: []string{"rb"},
		Short:   "Registers a broker",
		Long:    `Registers a broker`,

		PreRunE: prepare(rbc, rbc.Context),
		RunE:    cmd.RunE(rbc),
	}
	result = rbc.addFlags(result)

	return result
}

func (rbc *RegisterBrokerCmd) addFlags(command *cobra.Command) *cobra.Command {
	command.Flags().StringVarP(&rbc.basicString, "basic", "b", "", "Sets the username and password for basic authentication. Format is <username:password>.")
	cmd.AddFormatFlag(command.Flags())

	return command
}

// Validate validates command's arguments
func (rbc *RegisterBrokerCmd) Validate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Name and URL are required")
	}

	if rbc.basicString == "" {
		return fmt.Errorf("--basic flag is required")
	}

	rbc.broker.Name = args[0]
	rbc.broker.URL = args[1]

	if len(args) > 2 {
		rbc.broker.Description = args[2]
	}

	return rbc.parseCredentials()
}

// Run runs the command's logic
func (rbc *RegisterBrokerCmd) Run() error {
	resultBroker, err := rbc.Client.RegisterBroker(&rbc.broker)
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(rbc.Output, rbc.outputFormat, resultBroker)
	output.Println(rbc.Output)
	return nil
}

// SetOutputFormat set output format
func (rbc *RegisterBrokerCmd) SetOutputFormat(format int) {
	rbc.outputFormat = format
}

// HideUsage hide command's usage
func (rbc *RegisterBrokerCmd) HideUsage() bool {
	return true
}

func (rbc *RegisterBrokerCmd) parseCredentials() error {
	if rbc.basicString != "" {
		splitBasicString := strings.Split(rbc.basicString, ":")
		if len(splitBasicString) != 2 {
			return errors.New("basic string is invalid")
		}
		user := splitBasicString[0]
		password := splitBasicString[1]
		basic := types.Basic{User: user, Password: password}
		rbc.broker.Credentials = &types.Credentials{Basic: basic}
	}

	return nil
}
