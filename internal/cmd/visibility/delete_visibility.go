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

package visibility

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/errors"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"
)

// DeleteVisibilityCmd wraps the smctl delete-visibility command
type DeleteVisibilityCmd struct {
	*cmd.Context

	input io.Reader
	force bool

	id string
}

// NewDeleteVisibilityCmd returns new delete-visibility command with context
func NewDeleteVisibilityCmd(context *cmd.Context, input io.Reader) *DeleteVisibilityCmd {
	return &DeleteVisibilityCmd{Context: context, input: input}
}

// Validate validates command's arguments
func (dv *DeleteVisibilityCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id is required")
	}

	dv.id = args[0]

	return nil
}

// Run runs the command's logic
func (dv *DeleteVisibilityCmd) Run() error {
	dv.Parameters.FieldQuery = append(dv.Parameters.FieldQuery, fmt.Sprintf("id = %s", dv.id))

	err := dv.Client.DeleteVisibilities(&dv.Parameters)

	if respErr, ok := err.(errors.ResponseError); ok && respErr.StatusCode == http.StatusNotFound {
		output.PrintMessage(dv.Output, "Visibility not found.\n")
		return nil
	} else if err != nil {
		output.PrintMessage(dv.Output, "Could not delete visibility(s). Reason: ")
		return err
	}
	output.PrintMessage(dv.Output, "Visibility successfully deleted.\n")

	return nil
}

// HideUsage hide command's usage
func (dv *DeleteVisibilityCmd) HideUsage() bool {
	return true
}

// AskForConfirmation asks the user to confirm deletion
func (dv *DeleteVisibilityCmd) AskForConfirmation() (bool, error) {
	if !dv.force {
		message := fmt.Sprintf("Do you really want to delete visibilities with ids [%s] (Y/n): ", dv.id)
		return cmd.CommonConfirmationPrompt(message, dv.Context, dv.input)
	}
	return true, nil
}

// PrintDeclineMessage prints confirmation decline message to the user
func (dv *DeleteVisibilityCmd) PrintDeclineMessage() {
	cmd.CommonPrintDeclineMessage(dv.Output)
}

// Prepare returns cobra command
func (dv *DeleteVisibilityCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "delete-visibility [id]",
		Aliases: []string{"dv"},
		Short:   "Deletes visibility",
		Long:    `Delete one or more visibilities by name.`,
		PreRunE: prepare(dv, dv.Context),
		RunE:    cmd.RunE(dv),
	}

	result.Flags().BoolVarP(&dv.force, "force", "f", false, "Force delete without confirmation")
	cmd.AddCommonQueryFlag(result.Flags(), &dv.Parameters)

	return result
}
