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
	"strings"

	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/spf13/cobra"
)

// GetInstanceCmd wraps the smctl list-brokers command
type GetInstanceCmd struct {
	*cmd.Context

	instanceName   string
	outputFormat   output.Format
	instanceParams *bool
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
		GeneralParams: gb.Parameters.GeneralParams,
	})
	if err != nil {
		return err
	}
	if len(instances.ServiceInstances) < 1 {
		output.PrintMessage(gb.Output, "No instance found with name: %s", gb.instanceName)
		return nil
	}
	if *gb.instanceParams {
		return gb.printParameters(instances)
	}

	resultInstances := &types.ServiceInstances{Vertical: true}
	for _, instance := range instances.ServiceInstances {
		inst, err := gb.Client.GetInstanceByID(instance.ID, &gb.Parameters)
		if err != nil {
			// The instance could be deleted after List and before Get
			if strings.Contains(err.Error(), "StatusCode: 404") {
				continue
			}
			return err
		}
		resultInstances.ServiceInstances = append(resultInstances.ServiceInstances, *inst)
	}

	if len(resultInstances.ServiceInstances) < 1 {
		output.PrintMessage(gb.Output, "No instance found with name: %s", gb.instanceName)
		return nil
	}

	output.PrintServiceManagerObject(gb.Output, gb.outputFormat, resultInstances)
	output.Println(gb.Output)
	return nil
}

func (gb *GetInstanceCmd) printParameters(instances *types.ServiceInstances) error {

	for _, instance := range instances.ServiceInstances {
		parameters, err := gb.Client.GetInstanceParameters(instance.ID, &gb.Parameters)
		if err != nil {
			// The instance could be deleted after List and before Get
			if strings.Contains(err.Error(), "StatusCode: 404") {
				continue
			}
			output.PrintMessage(gb.Output, "Unable to show configuration parameters for service instance id: %s\n", instance.ID)
			output.PrintMessage(gb.Output, "The error: %s\n\n", err)
			continue
		}
		if len(parameters) == 0 {
			output.PrintMessage(gb.Output, "No configuration parameters are set for service instance id: %s\n\n", instance.ID)
			continue
		}
		output.PrintMessage(gb.Output, "Showing configuration parameters for service instance id: %s \n", instance.ID)
		output.PrintMessage(gb.Output, "The parameters: \n")

		output.PrintMessage(gb.Output, "%s \n\n", output.PrintParameters(parameters))
	}

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
		Aliases: []string{"gi"},
		Short:   "Get single instance",
		Long:    `Get single instance by its name`,
		PreRunE: prepare(gb, gb.Context),
		RunE:    cmd.RunE(gb),
	}
	gb.instanceParams = result.PersistentFlags().Bool("show-instance-params", false, "Show the service instance configuration parameters")
	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &gb.Parameters)

	return result
}
