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
	platformName string
	prepare      cmd.PrepareFunc
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

	instanceID := instances.ServiceInstances[0].ID
	if len(instances.ServiceInstances) > 1 && gb.platformName == "" {
		output.PrintMessage(gb.Output, "More than 1 instance found with this name. Use --platform flag to provide differentiator")
		return nil
	}

	if len(instances.ServiceInstances) > 1 && gb.platformName != "" {
		instanceID, err = gb.getInstanceIDByPlatformName(instances)
		if err != nil {
			output.PrintMessage(gb.Output, err.Error())
			return nil
		}
	}

	instance, err := gb.Client.GetInstanceByID(instanceID, &query.Parameters{
		GeneralParams: []string{
			"last_op=true",
		},
	})
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(gb.Output, gb.outputFormat, instance)
	output.Println(gb.Output)

	return nil
}

func (gb *GetInstanceCmd) getInstanceIDByPlatformName(instances *types.ServiceInstances) (string, error) {
	for _, instance := range instances.ServiceInstances {
		platforms, err := gb.Client.ListPlatforms(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("name eq '%s'", gb.platformName),
			},
		})
		if err != nil {
			return "", err
		}
		if len(platforms.Platforms) < 1 {
			return "", fmt.Errorf("No platform found with name %s", gb.platformName)
		}

		if instance.PlatformID == platforms.Platforms[0].ID {
			return instance.ID, nil
		}
	}

	return "", fmt.Errorf("No matching instance name %s and platform name %s", gb.instanceName, gb.platformName)
}

// Validate validates command's arguments
func (gb *GetInstanceCmd) Validate(args []string) error {
	if len(args) < 1 {
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
	gb.prepare = prepare
	result := &cobra.Command{
		Use:     "get-instance [name]",
		Aliases: []string{"gb"},
		Short:   "Get single instance",
		Long:    `Get single instance by its name`,
		PreRunE: gb.prepare(gb, gb.Context),
		RunE:    cmd.RunE(gb),
	}

	result.Flags().StringVarP(&gb.platformName, "platform", "p", "", "platform name as a differentiator between instances with same names")

	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &gb.Parameters)

	return result
}
