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

package instance

import (
	"fmt"

	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/spf13/cobra"
)

// GetInstanceCmd wraps the smctl list-brokers command
type GetInstanceCmd struct {
	*cmd.Context

	instanceName string
	outputFormat output.Format
}

// NewGetInstanceCmd returns new get status command with context
func NewGetInstanceCmd(context *cmd.Context) *GetInstanceCmd {
	return &GetInstanceCmd{Context: context}
}

// Run runs the command's logic
func (gb *GetInstanceCmd) Run() error {
	instances, err := gb.Client.ListInstances(&query.Parameters{
		FieldQuery: []string{
			fmt.Sprintf("name eq '%s'", gb.instanceName),
		},
	})
	if err != nil {
		return err
	}
	if len(instances.ServiceInstances) < 1 {
		output.PrintMessage(gb.Output, "No instance found with name: %s", gb.instanceName)
		return nil
	}

	resultInstances := &types.ServiceInstances{}
	for _, instance := range instances.ServiceInstances {
		inst, err := gb.Client.GetInstanceByID(instance.ID, &query.Parameters{
			GeneralParams: []string{
				"last_op=true",
			},
		})
		if err != nil {
			return err
		}
		resultInstances.ServiceInstances = append(resultInstances.ServiceInstances, *inst)
	}

	output.PrintServiceManagerObject(gb.Output, gb.outputFormat, resultInstances)
	output.Println(gb.Output)

	return nil
}

// Validate validates command's arguments
func (gb *GetInstanceCmd) Validate(args []string) error {
	if len(args) < 1 || len(args[0]) < 1 {
		return fmt.Errorf("instance name is required")
	}

	gb.instanceName = args[0]

	return nil
}

// SetOutputFormat set output format
func (gb *GetInstanceCmd) SetOutputFormat(format output.Format) {
	gb.outputFormat = format
}

// HideUsage hide command's usage
func (gb *GetInstanceCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (gb *GetInstanceCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "get-instance [name]",
		Aliases: []string{"gb"},
		Short:   "Get single instance",
		Long:    `Get single instance by its name`,
		PreRunE: prepare(gb, gb.Context),
		RunE:    cmd.RunE(gb),
	}

	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &gb.Parameters)

	return result
}
