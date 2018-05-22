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

	"github.com/Peripli/service-manager-cli/pkg/smclient"

	"encoding/json"
	"errors"

	"github.com/spf13/cobra"
)

// RegisterBrokerCmd wraps the smctl register-broker command
type RegisterBrokerCmd struct {
	*cmd.Context

	broker types.Broker

	credentialsJSON string
	basicString     string

	outputFormat int
}

// NewRegisterBrokerCmd returns new register-broker command with context
func NewRegisterBrokerCmd(context *cmd.Context) *RegisterBrokerCmd {
	return &RegisterBrokerCmd{Context: context, broker: types.Broker{}}
}

func (rbc *RegisterBrokerCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "register-broker [name] [url] <description>",
		Aliases: []string{"rb"},
		Short:   "Registers a broker",
		Long:    `Registers a broker`,

		PreRunE: cmd.PreRunE(rbc, rbc.Context),
		RunE:    cmd.RunE(rbc),
	}
}

func (rbc *RegisterBrokerCmd) addFlags(command *cobra.Command) *cobra.Command {
	command.Flags().StringVarP(&rbc.credentialsJSON, "credentials", "c", "", "Sets the authentication type and credentials with a json string. Format is <'json-string'>.")
	command.Flags().StringVarP(&rbc.basicString, "basic", "b", "", "Sets the username and password for basic authentication. Format is <username:password>.")
	cmd.AddFormatFlag(command.Flags())

	return command
}

// Validate validates command's arguments
func (rbc *RegisterBrokerCmd) Validate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("requires at least 2 args")
	}

	if rbc.credentialsJSON == "" && rbc.basicString == "" {
		return fmt.Errorf("requires either --credentials or --basic flag")
	}

	if rbc.credentialsJSON != "" && rbc.basicString != "" {
		return fmt.Errorf("duplicate credentials declaration with --credentials and --basic flags")
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

// SetSMClient set the SM client
func (rbc *RegisterBrokerCmd) SetSMClient(client smclient.Client) {
	rbc.Client = client
}

// SetOutputFormat set output format
func (rbc *RegisterBrokerCmd) SetOutputFormat(format int) {
	rbc.outputFormat = format
}

// HideUsage hide command's usage
func (rbc *RegisterBrokerCmd) HideUsage() bool {
	return true
}

// Command returns cobra command
func (rbc *RegisterBrokerCmd) Command() *cobra.Command {
	result := rbc.buildCommand()
	result = rbc.addFlags(result)

	return result
}

func (rbc *RegisterBrokerCmd) parseCredentials() error {
	if rbc.credentialsJSON != "" {
		credentials := &types.Credentials{}
		if err := json.Unmarshal([]byte(rbc.credentialsJSON), &credentials); err != nil {
			return errors.New("credentials string is invalid")
		}
		rbc.broker.Credentials = credentials
	}

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
