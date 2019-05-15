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
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

// DeleteVisibilityCmd wraps the smctl delete-visibility command
type DeleteVisibilityCmd struct {
	*cmd.Context

	input io.Reader
	force bool

	ids []string
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

	dv.ids = args

	return nil
}

// Run runs the command's logic
func (dv *DeleteVisibilityCmd) Run() error {
	deletedVisibilities := make(map[string]bool)

	for _, toDelete := range dv.ids {
		err := dv.Client.DeleteVisibility(toDelete)
		if err != nil {
			output.PrintMessage(dv.Output, "Could not delete visibility %s. Reason %s\n", toDelete, err)
		} else {
			output.PrintMessage(dv.Output, "Visibility with id: %s successfully deleted\n", toDelete)
			deletedVisibilities[toDelete] = true
		}
	}

	for _, id := range dv.ids {
		if !deletedVisibilities[id] {
			output.PrintError(dv.Output, fmt.Errorf("visibility with id: %s was not found", id))
		}
	}

	return nil
}

// HideUsage hide command's usage
func (dv *DeleteVisibilityCmd) HideUsage() bool {
	return true
}

// AskForConfirmation asks the user to confirm deletion
func (dv *DeleteVisibilityCmd) AskForConfirmation() (bool, error) {
	if !dv.force {
		message := fmt.Sprintf("Do you really want to delete visibilities with ids [%s] (Y/n): ", strings.Join(dv.ids, ", "))
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
		Use:     "delete-visibility [id] <id2 <id3> ... <idN>>",
		Aliases: []string{"dv"},
		Short:   "Deletes visibility",
		Long:    `Delete one or more visibilities by name.`,
		PreRunE: prepare(dv, dv.Context),
		RunE:    cmd.RunE(dv),
	}

	result.Flags().BoolVarP(&dv.force, "force", "f", false, "Force delete without confirmation")

	return result
}
