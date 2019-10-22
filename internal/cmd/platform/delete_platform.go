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

package platform

import (
	"fmt"
	"io"
	"strings"

	"github.com/Peripli/service-manager-cli/internal/output"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

// DeletePlatformCmd wraps the smctl delete-platform command
type DeletePlatformCmd struct {
	*cmd.Context

	input io.Reader
	force bool

	name string
}

// NewDeletePlatformCmd returns new delete-platform command with context
func NewDeletePlatformCmd(context *cmd.Context, input io.Reader) *DeletePlatformCmd {
	return &DeletePlatformCmd{Context: context, input: input}
}

// Validate validates command's arguments
func (dpc *DeletePlatformCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("single [name] is required")
	}

	dpc.name = args[0]

	return nil
}

// Run runs the command's logic
func (dpc *DeletePlatformCmd) Run() error {
	dpc.Parameters.FieldQuery = append(dpc.Parameters.FieldQuery, fmt.Sprintf("name eq '%s'", dpc.name))

	if err := dpc.Client.DeletePlatforms(&dpc.Parameters); err != nil {
		if strings.Contains(err.Error(), "StatusCode: 404") {
			output.PrintMessage(dpc.Output, "Platform(s) not found.\n")
			return nil
		}
		output.PrintMessage(dpc.Output, "Could not delete platform(s). Reason: ")
		return err
	}
	output.PrintMessage(dpc.Output, "Platform(s) successfully deleted.\n")
	return nil
}

// HideUsage hide command's usage
func (dpc *DeletePlatformCmd) HideUsage() bool {
	return true
}

// AskForConfirmation asks the user to confirm deletion
func (dpc *DeletePlatformCmd) AskForConfirmation() (bool, error) {
	if !dpc.force {
		message := fmt.Sprintf("Do you really want to delete platforms with names [%s] (Y/n): ", dpc.name)
		return cmd.CommonConfirmationPrompt(message, dpc.Context, dpc.input)
	}
	return true, nil
}

// PrintDeclineMessage prints confirmation decline message to the user
func (dpc *DeletePlatformCmd) PrintDeclineMessage() {
	cmd.CommonPrintDeclineMessage(dpc.Output)
}

// Prepare returns cobra command
func (dpc *DeletePlatformCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "delete-platform [name]",
		Aliases: []string{"dp"},
		Short:   "Deletes platforms",
		Long:    `Delete one or more platforms with name.`,
		PreRunE: prepare(dpc, dpc.Context),
		RunE:    cmd.RunE(dpc),
	}

	result.Flags().BoolVarP(&dpc.force, "force", "f", false, "Force delete without confirmation")
	cmd.AddCommonQueryFlag(result.Flags(), &dpc.Parameters)

	return result
}
