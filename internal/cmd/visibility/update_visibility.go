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
	"encoding/json"
	"fmt"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/spf13/cobra"
)

// UpdateVisibilityCmd wraps smctl update-visibility command
type UpdateVisibilityCmd struct {
	*cmd.Context

	outputFormat      output.Format
	id                string
	updatedVisibility *types.Visibility
}

// NewUpdateVisibilityCmd returns new smctl update-visibility command with context
func NewUpdateVisibilityCmd(context *cmd.Context) *UpdateVisibilityCmd {
	return &UpdateVisibilityCmd{Context: context}
}

// Validate validates command's arguments
func (uv *UpdateVisibilityCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id is required")
	}

	uv.id = args[0]

	if len(args) < 2 {
		return fmt.Errorf("nothing to update. Visibility JSON is not provided")
	}

	if err := json.Unmarshal([]byte(args[1]), &uv.updatedVisibility); err != nil {
		return fmt.Errorf("visibility JSON is invalid. Reason: %s", err.Error())
	}

	return nil
}

// Run runs the command's logic
func (uv *UpdateVisibilityCmd) Run() error {
	updatedVisibility, err := uv.Client.UpdateVisibility(uv.id, uv.updatedVisibility)
	if err != nil {
		return err
	}
	output.PrintServiceManagerObject(uv.Output, uv.outputFormat, updatedVisibility)
	output.Println(uv.Output)
	return nil
}

// HideUsage hide command's usage
func (uv *UpdateVisibilityCmd) HideUsage() bool {
	return true
}

// SetOutputFormat set output format
func (uv *UpdateVisibilityCmd) SetOutputFormat(format output.Format) {
	uv.outputFormat = format
}

//Prepare returns cobra command
func (uv *UpdateVisibilityCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "update-visibility [id] <json-visibility>",
		Aliases: []string{"uv"},
		Short:   "Updates visibility",
		Long:    "Updates visibility by id.",
		PreRunE: prepare(uv, uv.Context),
		RunE:    cmd.RunE(uv),
	}

	cmd.AddFormatFlag(result.Flags())

	return result
}
