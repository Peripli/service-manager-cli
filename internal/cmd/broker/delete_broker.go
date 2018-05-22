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

	id           string
	outputFormat int
}

// NewDeleteBrokerCmd returns new list-brokers command with context
func NewDeleteBrokerCmd(context *cmd.Context) *DeleteBrokerCmd {
	return &DeleteBrokerCmd{Context: context}
}

func (dbc *DeleteBrokerCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "delete-broker [id]",
		Aliases: []string{"db"},
		Short:   "Deletes broker",
		Long:    `Delete broker with id.`,
		PreRunE: cmd.PreRunE(dbc, dbc.Context),
		RunE:    cmd.RunE(dbc),
	}
}

// Validate validates command's arguments
func (dbc *DeleteBrokerCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[id] is required")
	}

	dbc.id = args[0]

	return nil
}

// Run runs the command's logic
func (dbc *DeleteBrokerCmd) Run() error {
	if err := dbc.Client.DeleteBroker(dbc.id); err != nil {
		return err
	}

	output.PrintMessage(dbc.Output, "Broker with id: %s successfully deleted", dbc.id)
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
