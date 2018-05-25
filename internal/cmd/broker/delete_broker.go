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
	"fmt"

	"github.com/Peripli/service-manager-cli/internal/output"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
)

// DeleteBrokerCmd wraps the smctl list-brokers command
type DeleteBrokerCmd struct {
	*cmd.Context

	names []string
}

// NewDeleteBrokerCmd returns new list-brokers command with context
func NewDeleteBrokerCmd(context *cmd.Context) *DeleteBrokerCmd {
	return &DeleteBrokerCmd{Context: context}
}

func (dbc *DeleteBrokerCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "delete-broker [name] <name2 <name3> ... <nameN>>",
		Aliases: []string{"db"},
		Short:   "Deletes brokers",
		Long:    `Delete one or more brokers with name.`,
		PreRunE: cmd.PreRunE(dbc, dbc.Context),
		RunE:    cmd.RunE(dbc),
	}
}

// Validate validates command's arguments
func (dbc *DeleteBrokerCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[name] is required")
	}

	dbc.names = args

	return nil
}

// Run runs the command's logic
func (dbc *DeleteBrokerCmd) Run() error {
	allBrokers, err := dbc.Client.ListBrokers()
	if err != nil {
		return err
	}

	toDeleteBrokers := getBrokerByName(allBrokers, dbc.names)
	if len(toDeleteBrokers) < 1 {
		return fmt.Errorf("No brokers are found")
	}

	for _, toDelete := range toDeleteBrokers {
		if err := dbc.Client.DeleteBroker(toDelete.ID); err != nil {
			return err
		}
		output.PrintMessage(dbc.Output, "Broker with name: %s successfully deleted\n", toDelete.Name)
	}

	if len(toDeleteBrokers) < len(dbc.names) {
		output.PrintError(dbc.Output, fmt.Errorf("%d names were not found", len(dbc.names)-len(toDeleteBrokers)))
	}

	return nil
}

// SetSMClient set the SM client
func (dbc *DeleteBrokerCmd) SetSMClient(client smclient.Client) {
	dbc.Client = client
}

// HideUsage hide command's usage
func (dbc *DeleteBrokerCmd) HideUsage() bool {
	return true
}

// Command returns cobra command
func (dbc *DeleteBrokerCmd) Command() *cobra.Command {
	result := dbc.buildCommand()

	return result
}
