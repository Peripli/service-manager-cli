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
	"github.com/Peripli/service-manager-cli/pkg/errors"
	"io"
	"net/http"
	"strings"

	"github.com/Peripli/service-manager-cli/internal/util"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

// DeleteBrokerCmd wraps the smctl delete-broker command
type DeleteBrokerCmd struct {
	*cmd.Context
	prepare cmd.PrepareFunc

	input io.Reader
	force bool

	names []string
}

// NewDeleteBrokerCmd returns new delete-broker command with context
func NewDeleteBrokerCmd(context *cmd.Context, input io.Reader) *DeleteBrokerCmd {
	return &DeleteBrokerCmd{Context: context, input: input}
}

// Validate validates command's arguments
func (dbc *DeleteBrokerCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[name] is required")
	}

	for _, name := range args {
		if name != "" {
			dbc.names = append(dbc.names, name)
		}
	}

	return nil
}

// Run runs the command's logic
func (dbc *DeleteBrokerCmd) Run() error {
	fieldQuery := util.GetResourceByNamesQuery(dbc.names)
	err := dbc.Client.DeleteBrokersByFieldQuery(fieldQuery)
	if respErr, ok := err.(errors.ResponseError); ok && respErr.StatusCode == http.StatusNotFound{
		output.PrintMessage(dbc.Output, "Service Broker(s) not found.\n")
	} else if err != nil {
		output.PrintMessage(dbc.Output, "Could not delete broker(s). Reason: ")
		return err
	}
	output.PrintMessage(dbc.Output, "Service Broker(s) successfully deleted.\n")
	return nil
}

// HideUsage hide command's usage
func (dbc *DeleteBrokerCmd) HideUsage() bool {
	return true
}

// AskForConfirmation asks the user to confirm deletion
func (dbc *DeleteBrokerCmd) AskForConfirmation() (bool, error) {
	if !dbc.force {
		message := fmt.Sprintf("Do you really want to delete brokers with names [%s] (Y/n): ", strings.Join(dbc.names, ", "))
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
	dbc.prepare = prepare
	result := &cobra.Command{
		Use:     "delete-broker [name] <name2 <name3> ... <nameN>>",
		Aliases: []string{"db"},
		Short:   "Deletes brokers",
		Long:    `Delete one or more brokers with name.`,
		PreRunE: dbc.prepare(dbc, dbc.Context),
		RunE:    cmd.RunE(dbc),
	}

	result.Flags().BoolVarP(&dbc.force, "force", "f", false, "Force delete without confirmation")

	return result
}
