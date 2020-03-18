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
	"github.com/Peripli/service-manager-cli/pkg/query"
	"io"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

// DeleteBrokerCmd wraps the smctl delete-broker command
type DeleteBrokerCmd struct {
	*cmd.Context

	input io.Reader
	force bool

	name string
}

// NewDeleteBrokerCmd returns new delete-broker command with context
func NewDeleteBrokerCmd(context *cmd.Context, input io.Reader) *DeleteBrokerCmd {
	return &DeleteBrokerCmd{Context: context, input: input}
}

// Validate validates command's arguments
func (dbc *DeleteBrokerCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("single [name] is required")
	}

	dbc.name = args[0]

	return nil
}

// Run runs the command's logic
func (dbc *DeleteBrokerCmd) Run() error {
	toDeleteBrokers, err := dbc.Client.ListBrokers(&query.Parameters{
		FieldQuery: []string{
			fmt.Sprintf("name eq '%s'", dbc.name),
		},
	})
	if err != nil {
		return err
	}
	if len(toDeleteBrokers.Brokers) < 1 {
		output.PrintMessage(dbc.Output, "Service Broker not found.\n")
		return nil
	}
	location, err := dbc.Client.DeleteBroker(toDeleteBrokers.Brokers[0].ID, &dbc.Parameters)
	if err != nil {
		output.PrintMessage(dbc.Output, "Could not delete broker. Reason: ")
		return err
	}
	if len(location) != 0 {
		cmd.CommonHandleAsyncExecution(dbc.Context, location, fmt.Sprintf("Service Broker %s successfully scheduled for deletion. To see status of the operation use:\n", dbc.name))
		return nil
	}
	output.PrintMessage(dbc.Output, "Service Broker successfully deleted.\n")
	return nil
}

// HideUsage hide command's usage
func (dbc *DeleteBrokerCmd) HideUsage() bool {
	return true
}

// AskForConfirmation asks the user to confirm deletion
func (dbc *DeleteBrokerCmd) AskForConfirmation() (bool, error) {
	if !dbc.force {
		message := fmt.Sprintf("Do you really want to delete broker with name [%s] (Y/n): ", dbc.name)
		return cmd.CommonConfirmationPrompt(message, dbc.Context, dbc.input)
	}
	return true, nil
}

// PrintDeclineMessage prints confirmation decline message to the user
func (dbc *DeleteBrokerCmd) PrintDeclineMessage() {
	cmd.CommonPrintDeclineMessage(dbc.Output)
}

// Prepare returns cobra command
func (dbc *DeleteBrokerCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "delete-broker [name]",
		Aliases: []string{"db"},
		Short:   "Deletes broker",
		Long:    `Deletes broker by name.`,
		PreRunE: prepare(dbc, dbc.Context),
		RunE:    cmd.RunE(dbc),
	}

	result.Flags().BoolVarP(&dbc.force, "force", "f", false, "Force delete without confirmation")
	cmd.AddCommonQueryFlag(result.Flags(), &dbc.Parameters)
	cmd.AddModeFlag(result.Flags(), "sync")

	return result
}
