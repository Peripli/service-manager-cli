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
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"fmt"

	"github.com/spf13/cobra"
)

// TransferCmd wraps the smctl provision command
type TransferCmd struct {
	*cmd.Context

	instanceName   string
	instanceID     string
	fromPlatformID string
	toPlatformID   string

	outputFormat output.Format
}

// NewTransferCmd returns new transfer instance command with context
func NewTransferCmd(context *cmd.Context) *TransferCmd {
	return &TransferCmd{Context: context}
}

// Prepare returns cobra command
func (trc *TransferCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "transfer-instance [name]",
		Short: "Transfer instance in one platform to another in SM",
		Long:  `Transfer instance in one platform to another in SM`,

		PreRunE: prepare(trc, trc.Context),
		RunE:    cmd.RunE(trc),
	}

	result.Flags().StringVarP(&trc.instanceID, "id", "", "", "Id of the instance. Required in case when there are instances with same name")
	result.Flags().StringVarP(&trc.fromPlatformID, "from", "", "", "ID of the platform from which you want to move the instance")
	result.Flags().StringVarP(&trc.toPlatformID, "to", "", "", "ID of the platform to which you want to move the instance")
	cmd.AddFormatFlag(result.Flags())
	cmd.AddModeFlag(result.Flags(), "async")

	return result
}

// Validate validates command's arguments
func (trc *TransferCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("name is required")
	}
	trc.instanceName = args[0]

	if len(trc.fromPlatformID) == 0 {
		return fmt.Errorf("--from is required")
	}

	if len(trc.toPlatformID) == 0 {
		return fmt.Errorf("--to is required")
	}

	return nil
}

// Run runs the command's logic
func (trc *TransferCmd) Run() error {
	var instance *types.ServiceInstance
	instances, err := trc.Client.ListInstances(&query.Parameters{
		FieldQuery: []string{
			fmt.Sprintf("name eq '%s'", trc.instanceName),
			fmt.Sprintf("platform_id eq '%s'", trc.fromPlatformID),
		},
	})
	if err != nil {
		return err
	}
	if len(instances.ServiceInstances) == 0 {
		return fmt.Errorf("No instances found with name %s", trc.instanceName)
	}

	if len(instances.ServiceInstances) > 1 {
		if len(trc.instanceID) == 0 {
			return fmt.Errorf("More than 1 instance found with name %s. Use --id flag to specify one", trc.instanceName)
		}
		instance, err = trc.Client.GetInstanceByID(trc.instanceID, nil)
		if err != nil {
			return err
		}
	} else {
		instance = &instances.ServiceInstances[0]
	}

	resultInstance, location, err := trc.Client.UpdateInstance(instance.ID, &types.ServiceInstance{
		PlatformID: trc.toPlatformID,
	}, nil)
	if err != nil {
		return err
	}

	if len(location) != 0 {
		cmd.CommonHandleAsyncExecution(trc.Context, location, fmt.Sprintf("Service Instance %s successfully scheduled for transfer to platform with id %s. To see status of the operation use:\n", trc.instanceName, trc.toPlatformID))
		return nil
	}
	output.PrintServiceManagerObject(trc.Output, trc.outputFormat, resultInstance)
	output.Println(trc.Output)
	return nil
}

// SetOutputFormat set output format
func (trc *TransferCmd) SetOutputFormat(format output.Format) {
	trc.outputFormat = format
}

// HideUsage hide command's usage
func (trc *TransferCmd) HideUsage() bool {
	return true
}
