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
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/spf13/cobra"
)

// RegisterVisibilityCmd wraps the smctl register-visibility command
type RegisterVisibilityCmd struct {
	*cmd.Context

	visibility types.Visibility

	outputFormat output.Format
}

// NewRegisterVisibilityCmd returns new smctl register-visibility command with context
func NewRegisterVisibilityCmd(ctx *cmd.Context) *RegisterVisibilityCmd {
	return &RegisterVisibilityCmd{Context: ctx}
}

// Validate validates command's arguments
func (rv *RegisterVisibilityCmd) Validate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("platform-id and service-plan-id required but not provided")
	}
	rv.visibility.PlatformID = args[0]
	rv.visibility.ServicePlanID = args[1]
	return nil
}

// Run runs command's logic
func (rv *RegisterVisibilityCmd) Run() error {
	resultVisibility, err := rv.Client.RegisterVisibility(&rv.visibility, &rv.Parameters)
	if err != nil {
		return err
	}
	output.PrintServiceManagerObject(rv.Output, rv.outputFormat, resultVisibility)
	output.Println(rv.Output)
	return nil
}

// SetOutputFormat sets command's output format
func (rv *RegisterVisibilityCmd) SetOutputFormat(format output.Format) {
	rv.outputFormat = format
}

// HideUsage hide command's usage
func (rv *RegisterVisibilityCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (rv *RegisterVisibilityCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "register-visibility [platform id] [service plan id]",
		Aliases: []string{"rv"},
		Short:   "Registers a visibility",
		Long:    "Registers a visibility",

		PreRunE: prepare(rv, rv.Context),
		RunE:    cmd.RunE(rv),
	}

	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &rv.Parameters)

	return result
}
