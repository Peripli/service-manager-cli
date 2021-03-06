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
	"encoding/json"
	"fmt"

	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

// UpdatePlatformCmd wraps the smctl update-platform command
type UpdatePlatformCmd struct {
	*cmd.Context

	outputFormat          output.Format
	name                  string
	regenerateCredentials bool
	updatedPlatform       *types.Platform
}

// NewUpdatePlatformCmd returns new update-platform command with context
func NewUpdatePlatformCmd(context *cmd.Context) *UpdatePlatformCmd {
	return &UpdatePlatformCmd{Context: context}
}

// Validate validates command's arguments
func (upc *UpdatePlatformCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[name] is required")
	}

	upc.name = args[0]

	if len(args) < 2 && !upc.regenerateCredentials {
		return fmt.Errorf("nothing to update. Platform JSON is not provided")
	}

	if len(args) >= 2 {
		if err := json.Unmarshal([]byte(args[1]), &upc.updatedPlatform); err != nil {
			return fmt.Errorf("platform JSON is invalid. Reason: %s", err.Error())
		}
	}

	return nil
}

// Run runs the command's logic
func (upc *UpdatePlatformCmd) Run() error {
	upc.Parameters.FieldQuery = append(upc.Parameters.FieldQuery, fmt.Sprintf("name eq '%s'", upc.name))
	toUpdatePlatforms, err := upc.Client.ListPlatforms(&upc.Parameters)
	if err != nil {
		return err
	}
	if len(toUpdatePlatforms.Platforms) < 1 {
		return fmt.Errorf("platform with name %s not found", upc.name)
	}
	toUpdatePlatform := toUpdatePlatforms.Platforms[0]
	if upc.regenerateCredentials {
		upc.Parameters.GeneralParams = append(upc.Parameters.GeneralParams, "regenerateCredentials=true")
	}
	result, err := upc.Client.UpdatePlatform(toUpdatePlatform.ID, upc.updatedPlatform, &upc.Parameters)
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(upc.Output, upc.outputFormat, result)
	output.Println(upc.Output)

	return nil
}

// HideUsage hide command's usage
func (upc *UpdatePlatformCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (upc *UpdatePlatformCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "update-platform [name] <json_platform>",
		Aliases: []string{"up"},
		Short:   "Updates platform",
		Long: `Update platform with name.
Example:
smctl update-platform platform '{"name": "new-name", "description": "new-description", "type": "new-type"}'`,
		PreRunE: prepare(upc, upc.Context),
		RunE:    cmd.RunE(upc),
	}

	result.Flags().BoolVarP(&upc.regenerateCredentials, "regenerate-credentials", "c", false, "Whether to regenerate credentials")

	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &upc.Parameters)

	return result
}

// SetOutputFormat set output format
func (upc *UpdatePlatformCmd) SetOutputFormat(format output.Format) {
	upc.outputFormat = format
}
