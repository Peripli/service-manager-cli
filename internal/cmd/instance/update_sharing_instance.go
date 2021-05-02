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

type UpdateSharingCmd struct {
	*cmd.Context
	instanceName string
	instanceID   string
	outputFormat output.Format
	share        bool
	action       string
}

// NewUpdateSharingCmd returns new share/unshare instance command with context
func NewUpdateSharingCmd(context *cmd.Context, share bool) *UpdateSharingCmd {
	return &UpdateSharingCmd{Context: context, share: share}
}

// Prepare returns cobra command
func (shc *UpdateSharingCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		PreRunE: prepare(shc, shc.Context),
		RunE:    cmd.RunE(shc),
	}
	if shc.share {
		shc.action = "share"
		result.Use = "share-instance [name] --id service-instance-id "
		result.Short = "Share a service instance"
		result.Long = `Share a service instance so that it can be consumed from various platforms in your subaccount.
Instance can be shared only if it was created with the plan that supports instance sharing. For more information, see the documentation of the service whose instance you want to share.`
	} else {
		shc.action = "unshare"
		result.Use = "unshare-instance [name] --id service-instance-id "
		result.Short = "Unshare a service instance"
		result.Long = `Unshare a service instance to disable its consumption from any but the original platform in which it was created in your subaccount. If an instance you want to unshare has references, an error is returned`
	}
	result.Flags().StringVarP(&shc.instanceID, "id", "", "", cmd.INSTANCE_ID_DESCRIPTION)
	cmd.AddFormatFlag(result.Flags())
	cmd.AddModeFlag(result.Flags(), "async")

	return result
}

// Validate validates command's arguments
func (shc *UpdateSharingCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("service instance name is required")
	}
	shc.instanceName = args[0]
	return nil
}

// Run runs the command's logic
func (shc *UpdateSharingCmd) Run() error {
	if shc.instanceID == "" {
		instances, err := shc.Client.ListInstances(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("name eq '%s'", shc.instanceName),
			},
			GeneralParams: shc.Parameters.GeneralParams,
		})
		if err != nil {
			return err
		}
		if len(instances.ServiceInstances) == 0 {
			return fmt.Errorf(cmd.NO_INSTANCES_FOUND, shc.instanceName)
		}

		if len(instances.ServiceInstances) > 1 {
			return fmt.Errorf(cmd.FOUND_TOO_MANY_INSTANCES, shc.instanceName, shc.action)
		}

		shc.instanceID = instances.ServiceInstances[0].ID
	}
	resultInstance, location, err := shc.Client.UpdateInstance(shc.instanceID, &types.ServiceInstance{
		Shared: shc.share,
	}, nil)
	var message = ""
	if err != nil {
		output.PrintMessage(shc.Output, fmt.Sprintf("Couldn't %s the service instance. Reason:",shc.action))
		return err
	}

	if len(location) != 0 {
		if shc.share {
			message = fmt.Sprintf("Service instance \"%s\" successfully scheduled for sharing. To see the status of the operation, use: \n", shc.instanceName)

		} else {
			message = fmt.Sprintf("Service instance \"%s\" successfully scheduled for unsharing. To see the status of the operation, use: \n", shc.instanceName)
		}
		cmd.CommonHandleAsyncExecution(shc.Context, location, message)
		return nil
	}
	output.PrintServiceManagerObject(shc.Output, shc.outputFormat, resultInstance)
	output.Println(shc.Output)
	return nil
}

// SetOutputFormat set output format
func (shc *UpdateSharingCmd) SetOutputFormat(format output.Format) {
	shc.outputFormat = format
}
