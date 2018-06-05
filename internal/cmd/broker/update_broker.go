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
	"encoding/json"
	"fmt"

	"github.com/Peripli/service-manager-cli/internal/output"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

// UpdateBrokerCmd wraps the smctl update-broker command
type UpdateBrokerCmd struct {
	*cmd.Context

	outputFormat  int
	name          string
	updatedBroker *types.Broker
}

// NewUpdateBrokerCmd returns new update-broker command with context
func NewUpdateBrokerCmd(context *cmd.Context) *UpdateBrokerCmd {
	return &UpdateBrokerCmd{Context: context}
}

func (ubc *UpdateBrokerCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "update-broker [name] <json_broker>",
		Aliases: []string{"ub"},
		Short:   "Updates broker",
		Long: `Update broker with name.
Example:
smctl update-broker broker '{"name": "new-name", "description": "new-description", "broker-url": "http://broker.com", "credentials": { "basic": { "username": "admin", "password": "admin" } }}'`,
		PreRunE: cmd.PreRunE(ubc, ubc.Context),
		RunE:    cmd.RunE(ubc),
	}
}

// Validate validates command's arguments
func (ubc *UpdateBrokerCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[name] is required")
	}

	ubc.name = args[0]

	if len(args) < 2 {
		return fmt.Errorf("Nothing to update. Broker JSON is not provided")
	}

	if err := json.Unmarshal([]byte(args[1]), &ubc.updatedBroker); err != nil {
		return fmt.Errorf("broker JSON is invalid. Reason: %s", err.Error())
	}

	return nil
}

// Run runs the command's logic
func (ubc *UpdateBrokerCmd) Run() error {
	brokers, err := ubc.Client.ListBrokers([]string{ubc.name})
	if err != nil {
		return err
	}

	if len(brokers.Brokers) < 1 {
		return fmt.Errorf("broker with name %s not found", ubc.name)
	}

	toUpdateBroker := brokers.Brokers[0]
	result, err := ubc.Client.UpdateBroker(toUpdateBroker.ID, ubc.updatedBroker)
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(ubc.Output, ubc.outputFormat, result)
	output.Println(ubc.Output)

	return nil
}

// SetSMClient set the SM client
func (ubc *UpdateBrokerCmd) SetSMClient(client smclient.Client) {
	ubc.Client = client
}

// HideUsage hide command's usage
func (ubc *UpdateBrokerCmd) HideUsage() bool {
	return true
}

// Command returns cobra command
func (ubc *UpdateBrokerCmd) Command() *cobra.Command {
	result := ubc.buildCommand()
	result = ubc.addFlags(result)

	return result
}

func (ubc *UpdateBrokerCmd) addFlags(command *cobra.Command) *cobra.Command {
	cmd.AddFormatFlag(command.Flags())
	return command
}

// SetOutputFormat set output format
func (ubc *UpdateBrokerCmd) SetOutputFormat(format int) {
	ubc.outputFormat = format
}
